package repositories

import (
	"context"
	"sync"

	"github.com/gordon-zhiyong/beehive-api/internal/model"
	"github.com/gordon-zhiyong/beehive-api/pkg/snowflake"
	"gorm.io/gorm"
)

// Task task repository
type Task struct {
}

var task *Task
var tOnce sync.Once
var defaultStatus uint8 = 1

// NewTask return task repository instance
func NewTask() *Task {
	tOnce.Do(func() {
		task = &Task{}
	})
	return task
}

func (t *Task) Create(ctx context.Context, m *model.Task) {
	db := ctx.Value("db").(*gorm.DB)
	t := &model.Task{
		ID:      snowflake.Generate(),
		EventID: m.ID,
		Topic:   m.Topic,
		Payload: m.Payload,
		Status:  defaultStatus,
	}
	db.Create(t)
}
