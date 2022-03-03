package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gordon-zhiyong/beehive-api/internal/model"
	"github.com/gordon-zhiyong/beehive-api/internal/repositories"
	"github.com/gordon-zhiyong/beehive-api/pkg/capsule"
)

// Topic topic service
type Topic struct {
	repo *repositories.Topic
}

// ErrTopicAlreadyExists topic already exists error
var ErrTopicAlreadyExists = repositories.ErrTopicAlreadyExists

// ErrSystemInternal system internal error
var ErrSystemInternal = errors.New("system internal error")

var topicOnce sync.Once
var topic *Topic

// NewTopic create a new topic service instance
func NewTopic() *Topic {
	topicOnce.Do(func() {
		topic = &Topic{
			repo: repositories.NewTopic(),
		}
	})
	return topic
}

// Create create a new topic
func (t *Topic) Create(name string) error {
	ctx := context.WithValue(context.TODO(), "db", capsule.DB)
	err := t.repo.Create(ctx, name)
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrTopicAlreadyExists) {
		return ErrTopicAlreadyExists
	}

	capsule.Logger.WithField("name", name).Errorf("create topic fail. err:%s", err.Error())
	return ErrSystemInternal
}

// UpdateName update topic name
func (t *Topic) UpdateName(ID uint64, name string) error {
	ctx := context.WithValue(context.TODO(), "db", capsule.DB)
	err := t.repo.Update(ctx, ID, name)
	fmt.Println(err)

	if err == nil {
		return nil
	}

	if errors.Is(err, ErrTopicAlreadyExists) {
		return ErrTopicAlreadyExists
	}

	fields := map[string]interface{}{"name": name, "id": ID}
	capsule.Logger.WithFields(fields).Errorf("update topic fail. err:%s", err.Error())
	return ErrSystemInternal
}

// Delete delete a topic
func (t *Topic) Delete(ID uint64) error {
	ctx := context.WithValue(context.TODO(), "db", capsule.DB)
	err := t.repo.Delete(ctx, ID)
	if err == nil {
		return nil
	}

	capsule.Logger.WithField("id", ID).Errorf("delete topic fail. err:%s", err.Error())
	return ErrSystemInternal
}

// Get get topic list
func (t *Topic) Get(filters ...Filter) []*model.Topic {
	ctx := context.WithValue(context.TODO(), "db", capsule.DB)
	return t.repo.Get(ctx, filters...)
}
