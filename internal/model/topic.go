package model

import "time"

// Topic topic model
type Topic struct {
	ID        uint64    `gorm:"column:id" db:"id" json:"id" form:"id"`                                 //数据表主键
	Name      string    `gorm:"column:name" db:"name" json:"name" form:"name"`                         //名称
	CreatedAt time.Time `gorm:"column:created_at" db:"created_at" json:"created_at" form:"created_at"` //创建时间
	UpdatedAt time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at" form:"updated_at"` //更新时间
}

// TableName return database table name
func (t Topic) TableName() string {
	return "topics"
}
