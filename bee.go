package beehive

import (
	"time"

	"github.com/pkg/errors"
)

type bee struct {
	task       chan func()
	recyleTime time.Time
	hive       *Beehive
}

func (b *bee) run() {
	b.hive.incRunning()
	go func() {
		defer func() {
			b.hive.decRunning()
			b.hive.workerCache.Put(b)

			if p := recover(); p != nil {
				if ph := b.hive.opts.PanicHandler; ph != nil {
					ph(p)
				} else {
					err := errors.Errorf("worker exits from panic: %v\n", p)
					b.hive.opts.Logger.Printf("worker exits from a panic: %v\n", err)
				}
			}
			b.hive.cond.Signal()
		}()

		for f := range b.task {
			if f == nil {
				return
			}

			f()

			if ok := b.hive.revertWorker(b); !ok {
				return
			}
		}
	}()
}
