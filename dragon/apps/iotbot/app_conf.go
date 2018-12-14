package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type AppConfig struct {
	GatewayId           string     `json:"gateway_id"`
	ServerUrl           string     `json:"server_url"`
	MqttBroker          MqttConfig `json:"mqtt_broker"`
	DeviceCycleInterval int        `json:"device_cycle_interval"`
}

func loadAppConfig(cfgPath string) (*AppConfig, error) {
	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		logger.Fatal("load app config file error, ", err.Error())
		return nil, err
	}

	fmt.Println(string(b))
	var cfg AppConfig
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		logger.Fatal("unmarshal app config file error, ", err.Error())
		return nil, err
	}

	fmt.Println(cfg.MqttBroker.Host, ":", cfg.MqttBroker.ClientId, ":", cfg.MqttBroker.UserName, ":", cfg.MqttBroker.Password)

	return &cfg, nil
}

func loadJsonConfig(cfgPath string, v interface{}) error {
	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		logger.Fatal("load config file ", cfgPath, " error, ", err.Error())
		return err
	}

	fmt.Println(string(b))
	err = json.Unmarshal(b, &v)
	if err != nil {
		logger.Fatal("unmarshal config file ", cfgPath, " error,", err.Error())
		return err
	}

	return nil
}
