package beehive

import (
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	OPEND = iota
	CLOSED
)

var (
	defaultLogger            = Logger(log.New(os.Stderr, "", log.LstdFlags))
	defaultSize              = 1024
	defaultInstance          = NewBeehive(defaultSize)
	defaultCleanIntervalTime = time.Minute
	defaultDumpFile          = "./beehive.json"
	ErrBeehiveClosed         = errors.New("this beehive has been closed")
	ErrTopicUnRegisted       = errors.New("this topic unregisted")
)

type Logger interface {
	Printf(format string, args ...interface{})
}

func Publish(topic string, data interface{}) {
	defaultInstance.Publish(topic, data)
}

func Release() {
	defaultInstance.Release()
}

func IsClosed() bool {
	return defaultInstance.IsClosed()
}

func Subscribe(topic string, listeners ...Handle) {
	defaultInstance.Subscribe(topic, listeners...)
}

func Cap() int {
	return defaultInstance.Cap()
}

func Running() int {
	return defaultInstance.Running()
}

func Reboot() {
	defaultInstance.Reboot()
}

func Free() {
	defaultInstance.Free()
}
