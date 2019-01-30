package utils

import (
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// NewLogger create logger tool to log
func NewLogger(path string) *logrus.Logger {

	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),

		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(time.Hour),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		//rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithMaxAge(time.Duration(7*24)*time.Hour),
	)

	hook := lfshook.NewHook(
		lfshook.WriterMap{
			logrus.FatalLevel: writer,
			logrus.ErrorLevel: writer,
			logrus.WarnLevel:  writer,
			logrus.DebugLevel: writer,
			logrus.InfoLevel:  writer,
		},
		&logrus.JSONFormatter{},
	)

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	//display call method
	//logger.SetReportCaller(true)
	logger.AddHook(hook)

	return logger
}
