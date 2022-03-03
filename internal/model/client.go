package model

import "time"

// Client 客户端
type Client struct {
	ID        uint64    `gorm:"column:id" db:"id" json:"id" form:"id"`                                 //数据表主键
	Name      string    `gorm:"column:name" db:"name" json:"name" form:"name"`                         //名称
	Secret    string    `gorm:"column:secret" db:"secret" json:"secret" form:"secret"`                 //秘钥
	Enabled   uint8     `gorm:"column:enabled" db:"enabled" json:"enabled" form:"enabled"`             //是否已激活(1:yes, 0:no)
	CreatedAt time.Time `gorm:"column:created_at" db:"created_at" json:"created_at" form:"created_at"` //创建时间
	UpdatedAt time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at" form:"updated_at"` //更新时间
}

// TableName return database table name
func (c Client) TableName() string {
	return "clients"
}
