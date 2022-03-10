package service

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gordon-zhiyong/beehive-api/pkg/capsule"
	"github.com/gordon-zhiyong/beehive-api/pkg/conf"
	"github.com/gordon-zhiyong/beehive-api/pkg/snowflake"
	bee "github.com/gordon-zhiyong/beehive-bee"
	"github.com/nsqio/go-nsq"
)

// Event event service
type Event struct {
	producer *nsq.Producer
}

var eventOnce sync.Once
var event *Event

// NewEvent create a event instance
func NewEvent() *Event {
	eventOnce.Do(func() {
		event = &Event{
			producer: newProducer(),
		}

	})

	return event
}

// newProducer create nsq producer
func newProducer() *nsq.Producer {
	address := conf.Conf.GetString("nsq.address")
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(address, config)
	if err != nil {
		capsule.Logger.Fatalf("nsq.NewProducer fail. err:%v", err)
	}
	return producer
}

// Publish publish event
func (e *Event) Publish(topic, payload, publisher string) error {
	event := &bee.Event{
		ID:          snowflake.Generate(),
		Topic:       topic,
		Publisher:   publisher,
		Payload:     payload,
		PublishedAt: time.Now(),
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = e.producer.Publish("NEW_EVENT", body)
	if err == nil {
		return nil
	}
	capsule.Logger.Errorf("deliver message to nsq fail. error:%v", err)
	return err
}
