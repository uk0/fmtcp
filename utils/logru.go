package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log *logrus.Logger

func LogInit(logLevel string) {
	log = logrus.New()
	//tmp, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.Out = os.Stdout
	level := logrus.DebugLevel
	log.SetLevel(logrus.ErrorLevel)
	switch {
	case logLevel == "debug":
		level = logrus.DebugLevel
	case logLevel == "info":
		level = logrus.InfoLevel
	case logLevel == "error":
		level = logrus.ErrorLevel
	default:
		level = logrus.DebugLevel
	}
	log.SetLevel(level)

}

func Info(v ...interface{}) {
	log.Info(v)
}

func Error(v ...interface{}) {
	log.Error(v)
}

func Debug(v ...interface{}) {
	log.Debug(v)
}
