package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gordon-zhiyong/beehive"
)

func listner(event *beehive.Event) error {
	fmt.Printf("listner(topic:%s, data:%v)\n", event.Topic, event.Data)
	time.Sleep(time.Millisecond)
	return nil
}

func listnerTwo(event *beehive.Event) error {
	fmt.Printf("listnerTwo(topic:%s, data:%v)\n", event.Topic, event.Data)
	time.Sleep(time.Millisecond)
	return nil
}

func producer(topic string) {
	t := time.NewTicker(time.Millisecond)
	defer t.Stop()

	for now := range t.C {
		if beehive.IsClosed() {
			return
		}
		beehive.Publish(topic, now.Format(time.RFC3339))
	}
}

func main() {
	beehive.Subscribe("testEvent", listner)
	beehive.Subscribe("testEvent2", listner, listnerTwo)
	beehive.Load()

	for i := 0; i < 1000; i++ {
		go func() {
			go producer("testEvent")
			go producer("testEvent")
			go producer("testEvent2")
			go producer("testEvent2")
		}()
	}

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	<-notify
	beehive.Release()
}
