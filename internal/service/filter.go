package service

import (
	"github.com/gordon-zhiyong/beehive-api/internal/repositories"
)

// Filter 过滤器
type Filter = repositories.Option

// WithName filter by name
func WithName(name string) Filter {
	return repositories.WithName(name)
}
