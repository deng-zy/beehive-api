package repositories

import (
	"context"
	"errors"
	"sync"

	"github.com/gordon-zhiyong/beehive-api/internal/model"
	"github.com/gordon-zhiyong/beehive-api/pkg/snowflake"
	"gorm.io/gorm"
)

// Topic topic repository
type Topic struct{}

var topicOnce sync.Once
var topic *Topic

// ErrTopicAlreadyExists topic already exists error
var ErrTopicAlreadyExists = errors.New("topic already exists")

// NewTopic create topic repository instance
func NewTopic() *Topic {
	topicOnce.Do(func() {
		topic = new(Topic)
	})
	return topic
}

// Create create a new topic
func (t *Topic) Create(ctx context.Context, name string) error {
	db, _ := ctx.Value("db").(*gorm.DB)

	record := &model.Topic{}
	db.First(record, "`name` = ?", name)
	if record.ID > 0 {
		return ErrTopicAlreadyExists
	}

	m := &model.Topic{
		ID:   snowflake.Generate(),
		Name: name,
	}

	return db.Create(m).Error
}

// Get get topic list
func (t *Topic) Get(ctx context.Context, filters ...Option) []*model.Topic {
	db, _ := ctx.Value("db").(*gorm.DB)

	var topics []*model.Topic
	db = db.Model(topics)

	for _, filter := range filters {
		filter(db)
	}

	db.Select("id", "name", "created_at", "updated_at").Find(&topics)
	return topics
}

// Update update topic name
func (t *Topic) Update(ctx context.Context, ID uint64, name string) error {
	db, _ := ctx.Value("db").(*gorm.DB)

	record := &model.Topic{}
	db.First(record, "`name` = ? and `id` != ?", name, ID)
	if record.ID > 0 {
		return ErrTopicAlreadyExists
	}

	return db.Model(record).Where("id=?", ID).Update("name", name).Error
}

// Delete delete a topic
func (t *Topic) Delete(ctx context.Context, ID uint64) error {
	db, _ := ctx.Value("db").(*gorm.DB)
	return db.Where("id=?", ID).Delete(&model.Topic{}).Error
}
