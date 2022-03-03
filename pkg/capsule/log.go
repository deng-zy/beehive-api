package capsule

import (
	"io"
	"os"
	"sync"

	"github.com/gordon-zhiyong/beehive-api/pkg/conf"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	loggers  = make(map[string]*logrus.Logger)
	logMutex sync.RWMutex
	config   *viper.Viper
)

func NewLogger(option ...string) *logrus.Logger {
	if config == nil {
		config = conf.Log
	}

	var name string
	if len(option) > 0 {
		name = option[0]
	} else {
		name = "beehive"
	}

	logMutex.RLock()
	logger, ok := loggers[name]
	logMutex.RUnlock()

	if ok {
		return logger
	}

	logMutex.Lock()
	defer logMutex.Unlock()
	logger = logrus.New()
	loggers[name] = logger

	writers := []io.Writer{os.Stdout}
	logFile := conf.Log.GetString(name + ".file")
	var (
		writer *os.File
		err    error
	)

	if writer, err = os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0755); err == nil {
		writers = append(writers, writer)
	}

	if formatter := config.GetString(name + ".formatter"); formatter == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	if level := config.GetUint32(name + ".level"); level != 0 {
		logger.SetLevel(logrus.Level(level))
	} else {
		logger.SetLevel(logrus.DebugLevel)
	}

	logger.SetLevel()

	logger.SetOutput(io.MultiWriter(writers...))
	return logger
}
