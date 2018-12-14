package main

import (
	"fmt"
	"log"

	"github.com/kardianos/service"
	"github.com/robfig/cron"
	"yonghui.cn/dragon/utils"
)

type BotService struct{}

func (svc *BotService) Start(s service.Service) error {
	log.Println("开始服务")
	go svc.run()
	return nil
}
func (svc *BotService) Stop(s service.Service) error {
	log.Println("停止服务")
	return nil
}
func (svc *BotService) run() {

	rootpath := utils.GetCurrentDirectory()
	appCfgPath := fmt.Sprintf("%s/config/app-cfg.json", rootpath)
	devCfgPath := fmt.Sprintf("%s/config/vlong-cfg.json", rootpath)
	ledCfgPath := fmt.Sprintf("%s/config/led-cfg.json", rootpath)

	// 这里放置程序要执行的代码……
	cfg, err := loadAppConfig(appCfgPath)
	if err != nil {
		logger.Fatal("load app config file error, ", err.Error())
		return
	}

	err = LoadDeviceConf(devCfgPath)
	if err != nil {
		logger.Fatal("load vlongsoft device config file error, ", err.Error())
		return
	}

	err = LoadLedConf(ledCfgPath)
	if err != nil {
		logger.Fatal("load led config file error, ", err.Error())
		return
	}

	initialRouter()
	go InitialMqtt(cfg.GatewayId, cfg.MqttBroker)

	go runTask(cfg.DeviceCycleInterval)
	//go StartCycleMonitor(cfg.DeviceCycleInterval)
}

func runTask(interval int) {
	task = cron.New()
	span := fmt.Sprintf("@every %ds", interval)
	task.AddFunc(span, func() {
		for _, dev := range terminals {
			fmt.Println("start collect")
			go collectState(dev)
		}
	})

	// task.AddFunc(span, func() {
	// 	for _, house := range greenHouses {
	// 		go syncLedDisplay(house)
	// 	}
	// })

	task.Start()
}
