package model

import "time"

// Event events table struct
type Event struct {
	ID          uint64    `gorm:"column:id" db:"id" json:"id" form:"id"`                                         //数据表主键
	Topic       string    `gorm:"column:topic" db:"topic" json:"topic" form:"topic"`                             //topic 主题
	Payload     string    `gorm:"column:payload" db:"payload" json:"payload" form:"payload"`                     //消息
	Publisher   string    `gorm:"column:publisher" db:"publisher" json:"publisher" form:"publisher"`             //发布者
	PublishedAt time.Time `gorm:"column:published_at" db:"published_at" json:"published_at" form:"published_at"` //发布时间 发布时间可能早于创建时间
	CreatedAt   time.Time `gorm:"column:created_at" db:"created_at" json:"created_at" form:"created_at"`         //创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at" form:"updated_at"`         //更新时间
}

// TableName return data table name
func (e *Event) TableName() string {
	return "events"
}
