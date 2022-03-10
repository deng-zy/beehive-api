package service

import (
	"context"
	"sync"
	"time"

	"github.com/gordon-zhiyong/beehive-api/internal/model"
	"github.com/gordon-zhiyong/beehive-api/internal/repositories"
	"github.com/gordon-zhiyong/beehive-api/pkg/capsule"
	"github.com/gordon-zhiyong/beehive-api/pkg/conf"
	"github.com/gordon-zhiyong/beehive-api/pkg/snowflake"
	"github.com/nsqio/go-nsq"
)

// Event event service
type Event struct {
	repo     *repositories.Event
	taskRepo *repositories.Task
	producer *nsq.Producer
}

var eventOnce sync.Once
var event *Event

// NewEvent create a event instance
func NewEvent() *Event {
	eventOnce.Do(func() {
		address := conf.App.GetString("nsq.address")
		config := nsq.NewConfig()
		producer, err := nsq.NewProducer(address, config)
		if err != nil {
			capsule.Logger.Fatalf("nsq.NewProducer fail. err:%v", err)
		}

		event = &Event{
			repo:     repositories.NewEvent(),
			taskRepo: repositories.NewTask(),
			producer: producer,
		}

	})

	return event
}

// Create create event
func (e *Event) Create(event *model.Event) error {
	ctx := context.WithValue(context.TODO(), "db", capsule.DB)
	if err := e.repo.Create(ctx, event); err != nil {
		return err
	}

	return e.taskRepo.Create(ctx, event)
}

// Publish publish event
func (e *Event) Publish(topic, payload, publisher string) error {
	event := &model.Event{
		ID:          snowflake.Generate(),
		Topic:       topic,
		Publisher:   publisher,
		Payload:     payload,
		PublishedAt: time.Now(),
	}
	if e.box == nil || !e.box.IsRunning() {
		return e.Create(event)
	}

	task := func() {
		e.Create(event)
	}

	e.box.Push(task)
	return nil
}
