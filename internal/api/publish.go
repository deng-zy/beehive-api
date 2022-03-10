package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gordon-zhiyong/beehive-api/internal/auth"
	"github.com/gordon-zhiyong/beehive-api/internal/request"
	"github.com/gordon-zhiyong/beehive-api/internal/service"
	"github.com/gordon-zhiyong/beehive-api/pkg/res"
)

// Publish event publish api
func Publish(c *gin.Context) {
	req := &request.Event{}
	err := c.ShouldBind(req)
	client, _ := c.MustGet("client").(*auth.ClientAuth)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, err))
		return
	}

	if !service.NewTopic().Exists(req.Topic) {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(TopicNotExists, ErrTopicNotExists))
		return
	}

	service.NewEvent().Publish(req.Topic, req.Payload, client.Name)
	c.JSON(http.StatusOK, res.JSONSuccess())
}

// MPublish public event with multi payload api
func MPublish(c *gin.Context) {
	req := &request.MPubReq{}
	err := c.ShouldBind(req)
	client, _ := c.MustGet("client").(*auth.ClientAuth)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, err))
		return
	}

	payloads := strings.Split(req.Payload, "\n")
	for _, payload := range payloads {
		if len(payload) > 2048 {
			c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(PayloadTooLarge, ErrPayloadTooLarge))
			return
		}
	}

	if !service.NewTopic().Exists(req.Topic) {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(TopicNotExists, ErrTopicNotExists))
		return
	}

	ser := service.NewEvent()
	for _, payload := range payloads {
		ser.Publish(req.Topic, payload, client.Name)
	}

	c.JSON(http.StatusOK, res.JSONSuccess())
}

// MPubWithMultiTopic public event with multi topic api
func MPubWithMultiTopic(c *gin.Context) {
	req := &request.MPubWithTopicReq{}
	err := c.ShouldBind(req)
	client, _ := c.MustGet("client").(*auth.ClientAuth)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, err))
		return
	}

	topics := strings.Split(req.Topic, "\n")
	if len(topics) > 100 {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(TopicTooLarge, ErrTopicTooLarge))
		return
	}

	for _, topic := range topics {
		if len(topic) > 64 {
			c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(TopicLengthTooLarge, ErrTopicLengthTooLarge))
			return
		}
	}

	if !service.NewTopic().ExistsWithMName(topics) {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(TopicNotExists, ErrTopicNotExists))
		return
	}

	ser := service.NewEvent()
	for _, topic := range topics {
		ser.Publish(topic, req.Payload, client.Name)
	}
	c.JSON(http.StatusOK, res.JSONSuccess())
}
