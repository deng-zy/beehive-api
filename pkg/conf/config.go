package conf

import (
	"errors"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var Conf *viper.Viper
var DB *viper.Viper
var Log *viper.Viper
var App *viper.Viper
var Auth *viper.Viper
var File string

var defaults map[string]string = map[string]string{
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
