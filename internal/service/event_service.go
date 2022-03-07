package service

import (
	"context"
	"sync"
	"time"

	"github.com/gordon-zhiyong/beehive-api/internal/app/box"
	"github.com/gordon-zhiyong/beehive-api/internal/model"
	"github.com/gordon-zhiyong/beehive-api/internal/repositories"
	"github.com/gordon-zhiyong/beehive-api/pkg/capsule"
	"github.com/gordon-zhiyong/beehive-api/pkg/snowflake"
)

// Event event service
type Event struct {
	repo *repositories.Event
	box  *box.Box
}

var eventOnce sync.Once
var event *Event

// NewEvent create a event instance
func NewEvent() *Event {
	eventOnce.Do(func() {
		event = &Event{
			repo: repositories.NewEvent(),
		}
	})

	return event
}

func (e *Event) SetBox(b *box.Box) {
	e.box = b
}

func (e *Event) Publish(topic, payload, publisher string) error {
	ID := snowflake.Generate()
	m := &model.Event{
		ID:          ID,
		Topic:       topic,
		Payload:     payload,
		Publisher:   publisher,
		PublishedAt: time.Now(),
	}

	ctx := context.WithValue(context.TODO(), "db", capsule.DB)
	err := e.repo.Create(ctx, m)
	if err != nil {
		return err
	}

	return nil
}
