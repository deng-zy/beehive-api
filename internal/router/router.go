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
	r.POST("/topics", api.CreateTopic)
	r.GET("/topics", api.GetTopics)
	r.PUT("/topics/:id", api.UpdateTopic)
	r.DELETE("/topics/:id", api.DeleteTopic)
	r.POST("/pub", api.Publish)
	r.POST("/mpub", api.MultiPublish)
}
