package service

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/gordon-zhiyong/beehive-api/internal/auth"
	"github.com/gordon-zhiyong/beehive-api/internal/model"
	"github.com/gordon-zhiyong/beehive-api/internal/repositories"
	"github.com/gordon-zhiyong/beehive-api/pkg/bytesconv"
	"github.com/gordon-zhiyong/beehive-api/pkg/capsule"
)

// Client Client service
type Client struct {
	repo *repositories.Client
}

var clientOnce sync.Once
var client *Client

// NewClient create a Client instance
func NewClient() *Client {
	clientOnce.Do(func() {
		client = &Client{
			repo: repositories.NewClient(),
		}
		rand.Seed(time.Now().UnixNano())
	})
	return client
}

// Create create a new client
func (c *Client) Create(name string) {
	ctx := context.WithValue(context.TODO(), "db", capsule.DB)
	secret := c.generateSecret()
	c.repo.Create(ctx, name, secret)
}

// Get get client list
func (c *Client) Get() []*model.Client {
	ctx := context.WithValue(context.TODO(), "db", capsule.DB)
	return c.repo.Get(ctx)
}

// Show show client info
func (c *Client) Show(ID uint64) *model.Client {
	ctx := context.WithValue(context.TODO(), "db", capsule.DB)
	return c.repo.WithId(ctx, ID)
}

// IssueToken issue auth token
func (c *Client) IssueToken(ID uint64, secret string) (string, error) {
	client := c.Show(ID)
	if client == nil {
		return "", errors.New("secret error")
	}

	if client.Secret != secret {
		return "", errors.New("secret error")
	}

	return auth.IssueToken(client.Id, client.Name)
}

func (c *Client) generateSecret() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	buff := make([]byte, 24)
	for i := range buff {
		buff[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return bytesconv.BytesToString(buff)
}
