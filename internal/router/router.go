package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gordon-zhiyong/beehive-api/internal/handler"
	"github.com/gordon-zhiyong/beehive-api/internal/middleware"
)

func Router(app *gin.Engine) {
	apiRouter(app.Group("/api"))
}

func apiRouter(r *gin.RouterGroup) {
	r.Use(middleware.ClientOAuth())
	r.POST("/client/oauth/token", handler.IssueClientToken)
	r.POST("/client/oauth/token/refresh", handler.RefreshClientToken)
	r.POST("/clients", handler.CreateClient)
	r.GET("/clients", handler.GetClients)
	r.GET("/client/:id", handler.ClientInfo)
	r.GET("/client", handler.ClientInfo)
}
