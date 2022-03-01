package main

import (
	"flag"

	"github.com/gordon-zhiyong/beehive-api/internal/app"
	"github.com/gordon-zhiyong/beehive-api/pkg/conf"
)

var configFile string
var env string

func init() {
	flag.StringVar(&configFile, "c", "./config/app.yaml", "configure file")
	flag.StringVar(&env, "e", "production", "application environment(local|test|production)")
}

func main() {
	flag.Parse()
	conf.Load(configFile, env)
	app.Run()
}
