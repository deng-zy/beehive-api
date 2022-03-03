package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gordon-zhiyong/beehive-api/internal/service"
	"github.com/gordon-zhiyong/beehive-api/pkg/res"
	"github.com/pkg/errors"
)

// CreateTopic create topic api
func CreateTopic(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	if name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, ErrInvalidParam))
		return
	}

	err := service.NewTopic().Create(name)
	if err == nil {
		c.JSON(http.StatusOK, res.JSONSuccess())
		return
	}

	if errors.Is(err, service.ErrTopicAlreadyExists) {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, err))
		return
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(SystemInternalError, err))
}

// UpdateTopic update topic name
func UpdateTopic(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if name == "" || id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, ErrInvalidParam))
		return
	}

	err := service.NewTopic().UpdateName(id, name)
	if err == nil {
		c.JSON(http.StatusOK, res.JSONSuccess())
		return
	}

	if errors.Is(err, service.ErrTopicAlreadyExists) {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, err))
		return
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(SystemInternalError, err))
}

// DeleteTopic delete topic with id
func DeleteTopic(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, ErrInvalidParam))
		return
	}

	service.NewTopic().Delete(id)
	c.JSON(http.StatusOK, res.JSONSuccess())
}

// GetTopics delete topic with id
func GetTopics(c *gin.Context) {
	topics := service.NewTopic().Get()
	c.JSON(http.StatusOK, res.JSONData(topics))
}
