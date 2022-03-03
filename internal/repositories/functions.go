package repositories

import "gorm.io/gorm"

// Option add query where
type Option func(*gorm.DB)

// WithID add query where id=xxxx
func WithID(ID uint64) Option {
	return func(db *gorm.DB) {
		db.Where("`id` = ?", ID)
	}
}

// WithName add query where name=xxxx
func WithName(name string) Option {
	return func(db *gorm.DB) {
		if name == "" {
			return
		}

		db.Where("`name` = ?", name)
	}
}
