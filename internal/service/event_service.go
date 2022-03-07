package service

import (
	"sync"
	"time"

	"github.com/gordon-zhiyong/beehive-api/internal/model"
)

// Event event service
type Event struct {
}

var eventOnce sync.Once
var event *Event

// NewEvent create a event instance
func NewEvent() *Event {
	eventOnce.Do(func() {
		event = &Event{}
	})

	return event
}

func (e *Event) Publish(topic, payload, publisher string) error {
	m := model.Event{
		Topic:       topic,
		Payload:     payload,
		Publisher:   publisher,
		PublishedAt: time.Now(),
	}

	return nil
}
