package repositories

import "gorm.io/gorm"

type Option func(*gorm.DB)

func WithId(id uint64) Option {
	return func(db *gorm.DB) {
		db.Where("`id` = ?", id)
	}
}
