# beehive
---
## Thanks

thanks https://github1s.com/panjf2000/ants

## Installl

```sh
go get github.com/gordon-zhiyong/beehive
```

## Usage

```go
package main

import (
    "fmt"

    "github.com/gordon-zhiyong/beehive"
)

func task(event *beehive.Event) error {
    fmt.Println("task say topic:", event.Topic, "data:", event.Data)
}

func taskTwo(event *beehive.Event) error {
    fmt.Println("taskTwo say topic:", event.Topic, "data:", event.Data)
}

func main() {
    beehive.Subscribe("demoEvent", task)
    beehive.Subscribe("muliListenerEvent", task, taskTwo)

    beehive.Publish("demoEvent")
    beehive.Publish("muliListenerEvent")

    time.Sleep(2 * time.Second)

    beehive.Release()
}

```