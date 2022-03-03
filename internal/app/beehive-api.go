package app

import (
	"github.com/gin-gonic/gin"
	"github.com/gordon-zhiyong/beehive-api/internal/router"
	"github.com/gordon-zhiyong/beehive-api/pkg/conf"
)

var envMap = map[string]string{
	"production": gin.ReleaseMode,
	"test":       gin.TestMode,
	"testing":    gin.TestMode,
	"local":      gin.DebugMode,
	"develop":    gin.DebugMode,
}

// Run bootstarap app
func Run() {
	env := conf.App.GetString("env")
	env, ok := envMap[env]
	if !ok {
		env = gin.ReleaseMode
	}
	gin.SetMode(env)

	app := gin.New()
	router.Router(app)
	app.Run(conf.App.GetString("host"))
}
