package capsule

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/gordon-zhiyong/beehive-api/pkg/conf"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	loggers  = make(map[string]*logrus.Logger)
	logMutex sync.RWMutex
	logConf  *viper.Viper
)

// NewLogger create log instance
func NewLogger(option ...string) *logrus.Logger {
	if logConf == nil {
		logConf = conf.Log
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

	if logFile != "" {
		writer, err = openLogFile(logFile)
		if err == nil {
			writers = append(writers, writer)
		}
	}

	if formatter := logConf.GetString(name + ".formatter"); formatter == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	logger.SetLevel(logrus.InfoLevel)
	if level := logConf.GetUint32(name + ".level"); level > 0 {
		l := logrus.Level(level)
		for _, v := range logrus.AllLevels {
			if v == l {
				logger.SetLevel(logrus.Level(level))
				break
			}
		}
	}

	logger.SetOutput(io.MultiWriter(writers...))
	return logger
}

func openLogFile(filename string) (file *os.File, err error) {
	dir := filepath.Dir(filename)
	pwd, _ := os.Getwd()
	path := filepath.Join(pwd, dir)
	if _, err := os.Stat(path); !os.IsExist(err) {
		os.Mkdir(path, 0755)
	}

	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	return
}
