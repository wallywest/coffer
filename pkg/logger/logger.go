package logger

import (
	"io"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	defaultLogger()
}

func defaultLogger() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{DisableColors: true}
	logger.Level = logrus.InfoLevel
	Logger = logger
}

func SetLogOutput(v io.Writer) {
	Logger.Out = v
}

func SetLogLevel(level string) {
	switch level {
	case "ERROR":
		Logger.Level = logrus.ErrorLevel
	case "DEBUG":
		Logger.Level = logrus.DebugLevel
	default:
		Logger.Level = logrus.InfoLevel
	}
}

func TestLogger() *logrus.Logger {
	defaultLogger()
	Logger.Out = ioutil.Discard
	return Logger
}
