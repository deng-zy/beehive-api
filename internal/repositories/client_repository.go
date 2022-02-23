package repositories

import (
	"context"
	"fmt"
	"sync"

	"github.com/gordon-zhiyong/beehive-api/internal/model"
	"github.com/gordon-zhiyong/beehive-api/pkg/snowflake"
	"gorm.io/gorm"
)

type Client struct {
}

var clientOnce sync.Once
var client *Client

func NewClient() *Client {
	clientOnce.Do(func() {
		client = &Client{}
	})
	return client
}

func (c *Client) Create(ctx context.Context, name, secret string) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	fmt.Println(ok)
	if !ok {
		return false
	}

	client := &model.Client{
		Id:      snowflake.Generate(),
		Name:    name,
		Secret:  secret,
		Enabled: 1,
	}

	db.Create(client)
	return true
}

func (c *Client) WithId(ctx context.Context, id uint64) *model.Client {
	clients := c.Get(ctx, WithId(id))
	if clients == nil || len(clients) < 1 {
		return nil
	}
	return clients[0]
}

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
