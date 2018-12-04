package main

import (
	"os"

	"github.com/kardianos/service"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"yonghui.cn/dragon/utils"
)

var logger *logrus.Entry
var task *cron.Cron

func main() {

	logger = utils.NewLogger("main entry", "logs/log")
	cfg := &service.Config{
		Name:        "iotbot",
		DisplayName: "YH-IoT-Gateway Service",
		Description: "YH-IoT-Gateway Service",
	}

	svc := &BotService{}
	s, err := service.New(svc, cfg)
	if err != nil {
		logger.Fatal("create system service error,", err.Error())
	}

	if len(os.Args) == 2 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			logger.Fatal("Execute control cmd error,", err.Error())
		}
	} else {
		err = s.Run()
		if err != nil {
			logger.Fatal("Execute control cmd error", err.Error())
		}
	}
}
