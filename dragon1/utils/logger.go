package utils

import (
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// NewLogger create logger tool to log
func NewLogger(name, path string) *logrus.Entry {
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		//	rotatelogs.WithLinkName(name),
		rotatelogs.WithMaxAge(time.Duration(3600)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(1)*time.Hour),
	)

	hook := lfshook.NewHook(
		lfshook.WriterMap{
			logrus.FatalLevel: writer,
			logrus.ErrorLevel: writer,
			logrus.WarnLevel:  writer,
			logrus.DebugLevel: os.Stdout,
			logrus.InfoLevel:  os.Stdout,
		},
		&logrus.JSONFormatter{},
	)

	logger := logrus.New()

	logger.AddHook(hook)
	return logger.WithFields(logrus.Fields{
		"module": name,
	})
}
