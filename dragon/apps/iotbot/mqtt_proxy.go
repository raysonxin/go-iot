package main

import (
	"fmt"
	"time"

	"yonghui.cn/dragon/drivers/mqtt"
)

type MqttConfig struct {
	Host     string `json:"host"`      // host
	ClientId string `json:"client_id"` // client id
	UserName string `json:"user_name"` // user name
	Password string `json:"password"`  // password
}

// mqttAdaptor mqtt adaptor instance
var mqttAdaptor *mqtt.Adaptor
var subRequestTopic string
var subResponseTopic string
var reportTopic string
var responseTopic string

func InitialMqtt(gatewayId string, conf MqttConfig) {
	subRequestTopic = fmt.Sprintf("gpgs/gateway-sub/request/%s", gatewayId)
	subResponseTopic = fmt.Sprintf("gpgs/gateway-sub/response/%s", gatewayId)
	reportTopic = fmt.Sprintf("gpgs/gateway/report/%s", gatewayId)
	responseTopic = fmt.Sprintf("gpgs/gateway/response/%s", gatewayId)

	ConnectMqtt(conf)
	SubscribeCommand()
}

// connect to mqtt broker
func ConnectMqtt(conf MqttConfig) {
	//	adaptor := mqtt.NewAdaptorWithAuth("tcp://10.0.90.55:1883", "dragon-test", "admin", "admin")
	mqttAdaptor = mqtt.NewAdaptorWithAuth(conf.Host, conf.ClientId, conf.UserName, conf.Password)
	err := mqttAdaptor.Connect()
	if err == nil {
		return
	}

	logger.Error("connect mqtt broker ", conf.Host, " failed,retry afert 5 second.")

	time.Sleep(5 * time.Second)
	for {
		err = mqttAdaptor.Connect()
		if err == nil {
			break
		}

		logger.Error("connect mqtt broker ", conf.Host, " failed,retry afert 5 second.")
		time.Sleep(5 * time.Second)
	}
	return
}

// subscribe command the app interest
func SubscribeCommand() {

	go func() {
		// subscribe remote control command
		ret := mqttAdaptor.Subscribe(subRequestTopic, processMqttMessage)
		if !ret {
			for {
				logger.Error("subscribe topic ", subRequestTopic, " failed,retry after 5 seconds.")
				time.Sleep(5 * time.Second)
				ret = mqttAdaptor.Subscribe(subRequestTopic, processMqttMessage)

				if ret {
					break
				}
			}

		}

		// subscribe response from another side
		ret = mqttAdaptor.Subscribe(subResponseTopic, processMqttMessage)
		if !ret {
			for {
				logger.Error("subscibe topic ", subResponseTopic, "failed,retry after 5 seconds.")
				time.Sleep(5 * time.Second)
				ret = mqttAdaptor.Subscribe(subResponseTopic, processMqttMessage)

				if ret {
					break
				}
			}
		}
	}()
}
