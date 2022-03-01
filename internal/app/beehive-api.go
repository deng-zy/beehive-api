package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gordon-zhiyong/beehive-api/internal/router"
	"github.com/gordon-zhiyong/beehive-api/pkg/conf"
)

func Run() {
	env := conf.App.GetString("env")
	fmt.Println("app", env)
	fmt.Println("app", conf.DB)
	envs := map[string]string{
		"production": gin.ReleaseMode,
		"test":       gin.TestMode,
		"testing":    gin.TestMode,
		"local":      gin.DebugMode,
		"develop":    gin.DebugMode,
	}

	env, ok := envs[env]
	if !ok {
		env = gin.ReleaseMode
	}
	fmt.Println(env)

	gin.SetMode(env)

	app := gin.New()
	router.Router(app)
	app.Run(conf.App.GetString("host"))
}
