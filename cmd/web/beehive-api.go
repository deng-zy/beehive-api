package main

import (
	"flag"

	"github.com/gordon-zhiyong/beehive-api/internal/app"
	"github.com/gordon-zhiyong/beehive-api/pkg/capsule"
	"github.com/gordon-zhiyong/beehive-api/pkg/conf"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "./config/app.yaml", "configure file")
}

func main() {
	flag.Parse()
	// 加载配置
	conf.Load(configFile)
	//初始各种连接和Logger
	capsule.Init()
	//运行
	app.Run()
}
