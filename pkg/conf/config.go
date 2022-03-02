package conf

import (
	"errors"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Conf 程序配置
var Conf *viper.Viper

// DB 数据库配置
var DB *viper.Viper

// Log 日志配置
var Log *viper.Viper

// App 应用配置
var App *viper.Viper

// Auth 认证相关配置
var Auth *viper.Viper

// File 配置文件
var File string

// defaults 默认配置
var defaults = map[string]string{
	"app.host":                  ":8080",
	"app.env":                   "production",
	"database.beehive.host":     "127.0.0.0.1",
	"database.beehive.port":     "3306",
	"database.beehive.user":     "root",
	"database.beehive.name":     "beehive",
	"database.beehive.password": "forget",
	"database.beehive.charset":  "utf8mb4",
	"database.maxIdle":          "10",
	"database.maxConns":         "30",
	"database.maxLifetime":      "10m",
	"database.default":          "beehive",
	"loggers.beehive.file":      "./logs/beehive.log",
	"loggers.beehive.formatter": "json",
	"loggers.beehive.level":     "Info",
}

func init() {
	Conf = viper.New()
}

// Load 从文件加载配置
func Load(file string) {
	File = file
	ext := path.Ext(File)
	dir := filepath.Dir(File)

	Conf.SetConfigType(strings.TrimLeft(ext, "."))
	Conf.AddConfigPath(dir)
	Conf.AutomaticEnv()
	Conf.SetConfigFile(file)

	err := Conf.ReadInConfig()
	if err != nil {
		panic(errors.New(err.Error()))
	}

	setDefault()
	setSub()
}

func setDefault() {
	for k, v := range defaults {
		Conf.SetDefault(k, v)
	}
}

func setSub() {
	DB = Conf.Sub("database")
	Log = Conf.Sub("loggers")
	Auth = Conf.Sub("oauth")
	App = Conf.Sub("app")
}
