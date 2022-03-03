package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gordon-zhiyong/beehive-api/internal/api"
	"github.com/gordon-zhiyong/beehive-api/internal/auth"
	"github.com/gordon-zhiyong/beehive-api/pkg/conf"
	"github.com/gordon-zhiyong/beehive-api/pkg/res"
)

type finder func(*gin.Context, string) string

var finders = map[string]finder{
	"header": func(c *gin.Context, key string) string {
		authHeader := c.Request.Header.Get(key)
		if authHeader == "" {
			return ""
		}

		if !strings.HasPrefix(authHeader, conf.Auth.GetString("token_name")) {
			return authHeader
		}

		parts := strings.SplitN(authHeader, " ", 2)
		return parts[1]
	},
	"query": func(c *gin.Context, key string) string {
		return c.Query(key)
	},
	"form": func(c *gin.Context, key string) string {
		return c.PostForm(key)
	},
	"cookie": func(c *gin.Context, name string) string {
		v, _ := c.Cookie(name)
		return v
	},
}

// ClientOAuth token验证中间件
func ClientOAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if shouldPassThrough(c.Request.URL.Path) {
			c.Next()
			return
		}

		var token string
		for _, method := range conf.Auth.GetStringSlice("lookup") {
			if len(token) > 0 {
				break
			}

			parts := strings.Split(method, ":")
			find, _ := finders[parts[0]]
			if find != nil {
				token = find(c, parts[1])
			}
		}

		if len(token) < 1 {
			c.AbortWithStatusJSON(http.StatusBadRequest, res.NewJSONError(api.InvalidParams, api.ErrTokenIsEmpty))
			return
		}

		client, err := auth.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, res.NewJSONError(api.AuthFail, api.ErrAuthFail))
			return
		}

		c.Set("client", client)
		c.Next()
	}
}

func shouldPassThrough(path string) bool {
	excepts := conf.Auth.GetStringSlice("excepts")
	path = strings.TrimLeft(path, "/")
	for _, except := range excepts {
		if except == path {
			return true
		}
	}

	return false
}
