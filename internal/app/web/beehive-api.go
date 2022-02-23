package web

import (
	"github.com/gin-gonic/gin"
	"github.com/gordon-zhiyong/beehive-api/internal/router"
)

var app = gin.New()

func Run() {
	router.Router(app)
	gin.Recovery()
	app.Run()
}
