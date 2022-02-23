package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gordon-zhiyong/beehive-api/internal/auth"
	"github.com/gordon-zhiyong/beehive-api/internal/handler"
	"github.com/gordon-zhiyong/beehive-api/pkg/res"
)

func ClientOAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if shouldPassThrough(c.Request.URL.Path) {
			c.Next()
			return
		}

		authHeader := c.Request.Header.Get(auth.TOKEN_LOOKUP)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, res.NewJsonError(handler.InvalidParams, handler.ErrInvalidParam))
			return
		}

		var token string
		if strings.HasPrefix(authHeader, auth.TOKEN_HEAD_NAME) {
			parts := strings.SplitN(authHeader, " ", 2)
			if parts[1] == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, res.NewJsonError(handler.InvalidParams, handler.ErrInvalidParam))
				return
			}
			token = parts[1]
		} else {
			token = authHeader
		}

		client, err := auth.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, res.NewJsonError(handler.AuthFail, handler.ErrAuthFail))
			return
		}

		c.Set("client", client)
		c.Next()
	}
}

func shouldPassThrough(path string) bool {
	excepts := []string{"api/client/oauth/token"}
	path = strings.TrimLeft(path, "/")
	for _, except := range excepts {
		if except == path {
			return true
		}
	}

	return false
}
