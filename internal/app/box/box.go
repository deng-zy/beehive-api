package box

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/gordon-zhiyong/beehive-api/internal/model"
	"github.com/gordon-zhiyong/beehive-api/internal/service"
	"github.com/pkg/errors"
)

type Box struct {
	workers int
	running int32
	exit    chan bool
	channel chan *model.Event
	ctx     context.Context
	service *service.Event
}

var boxOnce sync.Once
var box *Box

func NewBox(workers int, ctx context.Context) *Box {
	boxOnce.Do(func() {
		box = &Box{
			workers: workers,
			running: 0,
			ctx:     ctx,
			channel: make(chan *model.Event, 4096),
			service: service.NewEvent(),
		}
		service.SetBox(box)
	})
	return box
}

func (b *Box) Run() {
	if !atomic.CompareAndSwapInt32(&b.running, 0, 1) {
		return
	}

	for i := 0; i < b.workers; i++ {
		go b.work()
	}
}

func (b *Box) Push(e *model.Event) error {
	if !b.IsRunning() {
		return errors.New("service is shutting down")
	}

	go func() {
		b.channel <- e
	}()
	return nil
}

func (b *Box) Stop() {
	atomic.CompareAndSwapInt32(&b.running, 1, 0)
	<-b.exit
}

func (b *Box) IsRunning() bool {
	return b.running == 1
}

func (b *Box) work() {
	defer func() {
		b.exit <- true
	}()

	for b.IsRunning() {
		select {
		case <-b.ctx.Done():
			return
		case <-b.channel:
		}
	}
}
