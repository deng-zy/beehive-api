package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gordon-zhiyong/beehive-api/internal/request"
	"github.com/gordon-zhiyong/beehive-api/internal/service"
	"github.com/gordon-zhiyong/beehive-api/pkg/res"
)

// Publish event publish api
func Publish(c *gin.Context) {
	req := &request.Event{}
	err := c.ShouldBind(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, err))
		return
	}

	if !service.NewTopic().Exists(req.Topic) {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(TopicNotExists, ErrTopicNotExists))
		return
	}

	c.JSON(http.StatusOK, res.JSONSuccess())
}

// MultiPublish multi event publish api
func MultiPublish(c *gin.Context) {
	req := &request.Event{}
	err := c.ShouldBind(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, err))
		return
	}

	if !service.NewTopic().Exists(req.Topic) {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(TopicNotExists, ErrTopicNotExists))
		return
	}

	for _, msg := range strings.Split(req.Message, "\n") {
		m := &request.Event{
			Topic:   req.Topic,
			Message: msg,
		}
		fmt.Println(m)
	}

	c.JSON(http.StatusOK, res.JSONSuccess())
}
