package main

import (
	"yonghui.cn/dragon/drivers/mqtt"
)

type ControlService struct {
	mqttCli *mqtt.Adaptor
}

func NewControlService() *ControlService {
	return nil
}

func (p *ControlService) Start() error {

	return nil
}

func (p *ControlService) Stop() error {
	return nil
}
