package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gordon-zhiyong/beehive"
)

var instance = beehive.NewBeehive(2048, beehive.WithDumpFile("./example.json"), beehive.WithExpriyDuration(time.Minute))

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
		if instance.IsClosed() {
			return
		}
		instance.Publish(topic, now.Format(time.RFC3339))
	}
}

func main() {
	instance.Subscribe("testEvent", listner)
	instance.Subscribe("testEvent2", listner, listnerTwo)

	for i := 0; i < 1000; i++ {
		go func() {
			go producer("testEvent")
			go producer("testEvent")
			go producer("testEvent2")
			go producer("testEvent2")
		}()
	}

	go func() {
		t := time.After(time.Minute)
		<-t
		instance.Release()
		instance.Reboot()
		instance.Load()
	}()

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	<-notify
	instance.Release()
}
