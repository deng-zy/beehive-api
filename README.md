# beehive
一个高性能，分布式事件处理框架。

## Description
Go event system

## Installl

```sh
go get github.com/gordon-zhiyong/beehive
```

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/gordon-zhiyong/beehive"
)

func task(event *beehive.Event) error {
	fmt.Println("task say topic:", event.Topic, "data:", event.Data)
	return nil
}

func taskTwo(event *beehive.Event) error {
	fmt.Println("taskTwo say topic:", event.Topic, "data:", event.Data)
	return nil
}

func main() {
	beehive.Subscribe("demoEvent", task)
	beehive.Subscribe("muliListenerEvent", task, taskTwo)

	beehive.Publish("demoEvent", "hello")
	beehive.Publish("muliListenerEvent", "im event data")

	time.Sleep(2 * time.Second)

	beehive.Release()
}

```