package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gordon-zhiyong/beehive-api/internal/auth"
	"github.com/gordon-zhiyong/beehive-api/internal/request"
	"github.com/gordon-zhiyong/beehive-api/internal/service"
	"github.com/gordon-zhiyong/beehive-api/pkg/conf"
	"github.com/gordon-zhiyong/beehive-api/pkg/res"
)

// CreateClient 创建新的客户端
func CreateClient(c *gin.Context) {
	req := &request.Client{}
	err := c.ShouldBind(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, ErrInvalidParam))
		return
	}

	service.NewClient().Create(req.Name)
	c.JSON(http.StatusOK, res.JSONSuccess())
}

// GetClients 获取客户端列表
func GetClients(c *gin.Context) {
	c.JSON(http.StatusOK, res.JSONData(service.NewClient().Get()))
}

// ClientInfo 获取客户端信息
func ClientInfo(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		ID = c.DefaultQuery("id", "0")
	}

	clientID, _ := strconv.ParseUint(ID, 10, 64)
	if clientID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, ErrInvalidParam))
		return
	}

	client := service.NewClient().Show(clientID)
	if client == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, res.NewJSONError(ObjectNotFound, ErrObjectNotFound))
		return
	}
	c.JSON(http.StatusOK, res.JSONData(client))
}

// IssueClientToken 签发token
func IssueClientToken(c *gin.Context) {
	req := &auth.AuthRequest{}
	err := c.ShouldBind(req)
	if err != nil || req.ClientID < 1 || len(req.Secret) < 24 {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, ErrInvalidParam))
		return
	}

	var token string
	token, err = service.NewClient().IssueToken(req.ClientID, req.Secret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(AuthFail, err))
		return
	}

	c.JSON(http.StatusOK, res.JSONData(gin.H{
		"token":      token,
		"token_type": conf.Auth.GetString("token_name"),
		"expire_in":  conf.Auth.GetInt("expires"),
	}))
}

// RefreshClientToken 刷新客户端token
func RefreshClientToken(c *gin.Context) {
	client, ok := c.MustGet("client").(*auth.ClientAuth)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, ErrInvalidParam))
		return
	}

	token, err := auth.ReFreshToken(client)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(InvalidParams, ErrInvalidParam))
		return
	}

	c.JSON(http.StatusOK, res.JSONData(gin.H{
		"token":      token,
		"token_type": conf.Auth.GetString("token_name"),
		"expire_in":  conf.Auth.GetInt("expires"),
	}))
}
