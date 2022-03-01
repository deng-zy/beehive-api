package main

import (
	"flag"

	"github.com/gordon-zhiyong/beehive-api/internal/app"
	"github.com/gordon-zhiyong/beehive-api/pkg/conf"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "./config/app.yaml", "configure file")
}

func main() {
	flag.Parse()
	conf.Load(configFile)
	app.Run()
}
