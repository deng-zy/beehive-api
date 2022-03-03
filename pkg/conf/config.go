package conf

import (
	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gordon-zhiyong/beehive-api/pkg/bytesconv"
	"github.com/gordon-zhiyong/beehive-api/pkg/helper"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

var (
	// Conf 程序配置
	Conf *viper.Viper

	// DB 数据库配置
	DB *viper.Viper

	// Log 日志配置
	Log *viper.Viper

	// App 应用配置
	App *viper.Viper

	// Auth 认证相关配置
	Auth *viper.Viper

	// File 配置文件
	File string

	defaultFile = "./config/app.yaml"
	replacer    = strings.NewReplacer(".", "_")
)

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
	gotenv.Load()
	Conf = viper.New()
}

// Load 加载配置文件
func Load(file string) {
	File = file
	ext := path.Ext(File)
	dir := filepath.Dir(File)
	wd, _ := os.Getwd()

	Conf.SetConfigType(strings.TrimLeft(ext, "."))
	Conf.AddConfigPath(dir)
	Conf.AddConfigPath(wd)
	Conf.AutomaticEnv()
	Conf.SetConfigFile(file)

	content, err := helper.ReadFile(file)
	if err != nil {
		panic(err)
	}
	Conf.ReadConfig(replacePlaceholder(content))

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
	Auth = Conf.Sub("auth")
	App = Conf.Sub("app")

	DB.SetEnvKeyReplacer(replacer)
	Log.SetEnvKeyReplacer(replacer)
	Auth.SetEnvKeyReplacer(replacer)
	App.SetEnvKeyReplacer(replacer)
}

func replacePlaceholder(content []byte) io.Reader {
	tmp := bytesconv.BytesToString(content)
	for _, variable := range os.Environ() {
		parts := strings.Split(variable, "=")
		tmp = strings.ReplaceAll(tmp, "$"+parts[0], parts[1])
	}
	return bytes.NewReader(bytesconv.StringToBytes(tmp))
}
