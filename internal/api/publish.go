package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gordon-zhiyong/beehive-api/internal/request"
	"github.com/gordon-zhiyong/beehive-api/internal/service"
	"github.com/gordon-zhiyong/beehive-api/pkg/res"
)

// Publish publish message api
func Publish(c *gin.Context) {
	req := &request.Message{}
	err := c.ShouldBind(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, err))
		return
	}

	fmt.Println(req)
	if !service.NewTopic().Exists(req.Topic) {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(TopicNotExists, ErrTopicNotExists))
		return
	}

	c.JSON(http.StatusOK, res.JSONSuccess())
}

// MultiPublish multiPublish message api
func MultiPublish(c *gin.Context) {
	req := []*request.Message{}
	err := c.ShouldBind(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, err))
		return
	}

	fmt.Println(req)
	c.JSON(http.StatusOK, res.JSONSuccess())
}
