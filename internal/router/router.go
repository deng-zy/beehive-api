package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gordon-zhiyong/beehive-api/internal/api"
	"github.com/gordon-zhiyong/beehive-api/internal/middleware"
)

func Router(app *gin.Engine) {
	apiRouter(app.Group("/api"))
}

func apiRouter(r *gin.RouterGroup) {
	r.Use(middleware.ClientOAuth())
	r.POST("/client/oauth/token", api.IssueClientToken)
	r.POST("/client/oauth/token/refresh", api.RefreshClientToken)
	r.POST("/clients", api.CreateClient)
	r.GET("/clients", api.GetClients)
	r.GET("/client/:id", api.ClientInfo)
	r.GET("/client", api.ClientInfo)
}
