package repositories

import (
	"context"
	"fmt"
	"sync"

	"github.com/gordon-zhiyong/beehive-api/internal/model"
	"github.com/gordon-zhiyong/beehive-api/pkg/snowflake"
	"gorm.io/gorm"
)

// Client client repository
type Client struct {
}

var clientOnce sync.Once
var client *Client

// NewClient new client repository
func NewClient() *Client {
	clientOnce.Do(func() {
		client = &Client{}
	})
	return client
}

// Create create a new client
func (c *Client) Create(ctx context.Context, name, secret string) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	fmt.Println(ok)
	if !ok {
		return false
	}

	client := &model.Client{
		ID:      snowflake.Generate(),
		Name:    name,
		Secret:  secret,
		Enabled: 1,
	}

	db.Create(client)
	return true
}

// WithID get client info with client id
func (c *Client) WithID(ctx context.Context, id uint64) *model.Client {
	clients := c.Get(ctx, WithID(id))
	if clients == nil || len(clients) < 1 {
		return nil
	}
	return clients[0]
}

// Get get client list
func (c *Client) Get(ctx context.Context, options ...Option) []*model.Client {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil
	}

	var clients []*model.Client
	db = db.Model(&clients)
	for _, option := range options {
		option(db)
	}
	db.Select("id", "name", "secret", "enabled", "created_at", "updated_at").Find(&clients)
	return clients
}
