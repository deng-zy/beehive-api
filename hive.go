package beehive

import (
	"encoding/json"
	"io/ioutil"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var workerChnCap = func() int {
	if runtime.GOMAXPROCS(0) == 1 {
		return 0
	}
	return 1
}()

//Event 事件
type Event struct {
	Data  interface{} `json:"data"`
	Topic string      `json:"topic"`
}

//Beehive 事件
type Beehive struct {
	topics      map[string]string
	handles     map[string][]Handle
	bucket      chan *Event
	opts        *Options
	state       int32
	workerCache sync.Pool
	capacity    int32
	cond        *sync.Cond
	bees        *beeStack
	lock        sync.Locker
	rwLock      sync.RWMutex
	running     int32
	blockingNum int
}

// Handle 事件处理者
type Handle func(*Event) error

func loadOptions(options ...Option) (opts *Options) {
	opts = new(Options)
	for _, option := range options {
		option(opts)
	}

	return
}

//NewBeehive  returns a new Bus
func NewBeehive(size int, options ...Option) *Beehive {
	opts := loadOptions(options...)

	if opts.Logger == nil {
		opts.Logger = defaultLogger
	}

	if expiry := opts.ExpiryDuration; expiry < 1 {
		opts.ExpiryDuration = defaultCleanIntervalTime
	}

	if len(opts.DumpFile) == 0 || opts.DumpFile == "" {
		opts.DumpFile = defaultDumpFile
	}

	beehive := &Beehive{
		topics:   map[string]string{},
		handles:  map[string][]Handle{},
		opts:     opts,
		lock:     NewSpinLock(),
		capacity: int32(size),
	}

	beehive.bucket = beehive.newBucket()
	beehive.workerCache.New = func() interface{} {
		return &bee{
			task: make(chan func(), workerChnCap),
			hive: beehive,
		}
	}
	beehive.bees = newBeeStack(size)
	beehive.cond = sync.NewCond(beehive.lock)

	go beehive.purgePeriodically()
	go beehive.load()
	go beehive.pickHoney()

	return beehive
}

//Publish a new event
func (b *Beehive) Publish(topic string, data interface{}) error {
	if b.IsClosed() {
		return ErrBeehiveClosed
	}

	b.rwLock.RLock()
	_, exists := b.topics[topic]
	b.rwLock.RUnlock()

	if !exists {
		return ErrTopicUnRegisted
	}

	event := &Event{
		Topic: topic,
		Data:  data,
	}

	b.bucket <- event
	return nil
}

func (b *Beehive) Release() {
	atomic.StoreInt32(&b.state, CLOSED)
	b.lock.Lock()

	b.bees.reset()
	close(b.bucket)

	b.lock.Unlock()
	b.cond.Broadcast()

	b.dump()
}

func (b *Beehive) IsClosed() bool {
	return atomic.LoadInt32(&b.state) == CLOSED
}

//Subscribe a event
func (b *Beehive) Subscribe(topic string, listener ...Handle) {
	b.registerTopic(topic)
	b.registerHandle(topic, listener...)
}

func (b *Beehive) Cap() int {
	return int(atomic.LoadInt32(&b.capacity))
}

func (b *Beehive) Running() int {
	return int(atomic.LoadInt32(&b.running))
}

func (b *Beehive) Reboot() {
	if atomic.CompareAndSwapInt32(&b.state, CLOSED, OPEND) {
		b.bucket = b.newBucket()

		go b.purgePeriodically()
		go b.load()
		go b.pickHoney()
	}
}

func (b *Beehive) Free() int {
	c := b.Cap()
	if c < 0 {
		return -1
	}
	return c - b.Running()
}

func (b *Beehive) incRunning() {
	atomic.AddInt32(&b.running, 1)
}

func (b *Beehive) decRunning() {
	atomic.AddInt32(&b.capacity, -1)
}

func (b *Beehive) registerTopic(topic string) {
	b.rwLock.RLock()
	_, exists := b.topics[topic]
	b.rwLock.RUnlock()

	if exists {
		return
	}

	b.lock.Lock()
	b.topics[topic] = topic
	b.lock.Unlock()
}

func (b *Beehive) registerHandle(topic string, handle ...Handle) {
	b.rwLock.RLock()
	handles, exists := b.handles[topic]
	b.lock.Unlock()

	register := func(handle ...Handle) {
		b.lock.Lock()
		b.handles[topic] = handle
		b.lock.Unlock()
	}

	if exists {
		handles = append(handles, handle...)
		register(handles...)
		return
	}

	register(handle...)
}

func (b *Beehive) pickHoney() {
	wrap := func(event *Event, handle Handle) func() {
		return func() {
			handle(event)
		}
	}

	for !b.IsClosed() {
		event, ok := <-b.bucket
		if !ok {
			return
		}

		b.rwLock.RLock()
		handles, exists := b.handles[event.Topic]
		b.rwLock.RUnlock()

		if !exists {
			b.opts.Logger.Printf("topic:%s unregisted handler", event.Topic)
			continue
		}

		for _, handle := range handles {
			task := wrap(event, handle)
			worker := b.retriveWorker()
			if worker == nil {
				b.opts.Logger.Printf("topic:%s, get worker fail.", event.Topic)
				continue
			}
			worker.task <- task
		}
	}
}

func (b *Beehive) purgePeriodically() {
	heartbeat := time.NewTicker(b.opts.ExpiryDuration)
	defer heartbeat.Stop()

	for range heartbeat.C {
		if b.IsClosed() {
			break
		}

		b.lock.Lock()
		expiredWorkers := b.bees.retrieve(b.opts.ExpiryDuration)
		b.lock.Unlock()

		for i := range expiredWorkers {
			expiredWorkers[i].task <- nil
			expiredWorkers[i] = nil
		}

		if b.Running() == 0 {
			b.cond.Broadcast()
		}

	}
}

func (b *Beehive) retriveWorker() (w *bee) {
	spanWorker := func() {
		w = b.workerCache.Get().(*bee)
		w.run()
	}

	b.lock.Lock()

	w = b.bees.detach()
	if w != nil {
		b.lock.Unlock()
	} else if capacity := b.Cap(); capacity == -1 || capacity > b.Running() {
		b.lock.Unlock()
		spanWorker()
	} else {
		if b.opts.NonBlocking {
			b.lock.Unlock()
			return
		}
	retry:
		if b.opts.MaxBlockingTasks != 0 && b.blockingNum >= b.opts.MaxBlockingTasks {
			b.lock.Unlock()
			return
		}
		b.blockingNum++
		b.cond.Wait()
		var nw int

		if nw = b.Running(); nw == 0 {
			b.lock.Unlock()
			if !b.IsClosed() {
				spanWorker()
			}
			return
		}

		if w = b.bees.detach(); w == nil {
			if nw < capacity {
				b.lock.Unlock()
				spanWorker()
				return
			}
			goto retry
		}
		b.lock.Unlock()
	}
	return
}

func (b *Beehive) revertWorker(w *bee) bool {
	if capacity := b.Cap(); (capacity > 0 && b.Running() > capacity) || b.IsClosed() {
		return false
	}
	w.recyleTime = time.Now()
	b.lock.Lock()

	if b.IsClosed() {
		b.lock.Unlock()
		return false
	}

	err := b.bees.insert(w)
	if err != nil {
		b.lock.Unlock()
		return false
	}

	b.cond.Signal()
	b.lock.Unlock()

	return true
}

func (b *Beehive) newBucket() (bucket chan *Event) {
	bucket = make(chan *Event, b.capacity)
	return
}

func (b *Beehive) dump() {
	messages := []*Event{}
	for message := range b.bucket {
		messages = append(messages, message)
	}

	dump, err := json.Marshal(messages)
	if err != nil {
		b.opts.Logger.Printf("dump data fail:%v.", err)
		return
	}

	err = ioutil.WriteFile(b.opts.DumpFile, dump, 0644)
	if err != nil {
		b.opts.Logger.Printf("write to file fail:%v.", err)
	}
}

func (b *Beehive) load() {
	if len(b.opts.DumpFile) == 0 || b.opts.DumpFile == "" {
		return
	}

	contents, err := ioutil.ReadFile(b.opts.DumpFile)
	if err != nil {
		b.opts.Logger.Printf("open %s fail:%s", b.opts.DumpFile, err.Error())
		return
	}

	messages := []*Event{}

	err = json.Unmarshal(contents, &messages)
	if err != nil {
		b.opts.Logger.Printf("json.Unmarshal fail:%s", err.Error())
		return
	}

	for _, message := range messages {
		b.bucket <- message
	}
}
