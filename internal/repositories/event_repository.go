package repositories

import (
	"context"
	"sync"

	"github.com/gordon-zhiyong/beehive-api/internal/model"
	"gorm.io/gorm"
)

// Event event repository
type Event struct{}

var event *Event
var eOnce sync.Once

// NewEvent create a new event repository instance
func NewEvent() *Event {
	eOnce.Do(func() {
		event = &Event{}
	})
	return event
}

// Create insert into events table
func (e *Event) Create(ctx context.Context, m *model.Event) error {
	db, _ := ctx.Value("db").(*gorm.DB)
	return db.Create(m).Error
}
