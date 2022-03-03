package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gordon-zhiyong/beehive-api/internal/auth"
	"github.com/gordon-zhiyong/beehive-api/internal/service"
	"github.com/gordon-zhiyong/beehive-api/pkg/res"
)

// CreateClient 创建新的客户端
func CreateClient(c *gin.Context) {
	name := strings.Trim(c.PostForm("name"), " ")
	if len(name) < 1 || name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJsonError(InvalidParams, ErrInvalidParam))
		return
	}

	service.NewClient().Create(name)
	c.JSON(http.StatusOK, res.JsonSuccess())
}

// GetClients 获取客户端列表
func GetClients(c *gin.Context) {
	c.JSON(http.StatusOK, res.JsonData(service.NewClient().Get()))
}

func ClientInfo(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		ID = c.DefaultQuery("id", "0")
	}

	clientID, _ := strconv.ParseUint(ID, 10, 64)
	if clientID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJsonError(InvalidParams, ErrInvalidParam))
		return
	}

	client := service.NewClient().Show(clientID)
	if client == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, res.NewJsonError(ObjectNotFound, ErrObjectNotFound))
		return
	}
	c.JSON(http.StatusOK, res.JsonData(client))
}

func IssueClientToken(c *gin.Context) {
	req := &auth.AuthRequest{}
	err := c.ShouldBind(req)
	if err != nil || req.ClientID < 1 || len(req.Secret) < 24 {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJsonError(InvalidParams, ErrInvalidParam))
		return
	}

	var token string
	token, err = service.NewClient().IssueToken(req.ClientID, req.Secret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJsonError(AuthFail, err))
		return
	}

	c.JSON(http.StatusOK, res.JsonData(gin.H{
		"token":      token,
		"token_type": auth.TOKEN_HEAD_NAME,
		"expire_in":  auth.EXPIRES,
	}))
}

func RefreshClientToken(c *gin.Context) {
	client, ok := c.MustGet("client").(*auth.ClientAuth)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJsonError(InvalidParams, ErrInvalidParam))
		return
	}

	token, err := auth.ReFreshToken(client)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJsonError(InvalidParams, ErrInvalidParam))
		return
	}

	c.JSON(http.StatusOK, res.JsonData(gin.H{
		"token":      token,
		"token_type": auth.TOKEN_HEAD_NAME,
		"expire_in":  auth.EXPIRES,
	}))
}
