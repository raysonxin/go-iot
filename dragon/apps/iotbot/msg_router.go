package main

import (
	"encoding/json"
	"sync"
	"time"

	"yonghui.cn/dragon/drivers/mqtt"
	"yonghui.cn/dragon/protocol"
)

var reportChan chan protocol.Request
var responseChan chan protocol.Response
var seqLock sync.Mutex
var seq uint16

func initialRouter() {
	reportChan = make(chan protocol.Request, 0)
	responseChan = make(chan protocol.Response, 0)
	go routeMqttMessage()
}

// getSequenceNo counter
func getSequenceNo() uint16 {
	seqLock.Lock()
	defer seqLock.Unlock()
	seq++
	return seq
}

// process message from subscribe channel
func processMqttMessage(msg mqtt.Message) {

	switch msg.Topic() {
	case subRequestTopic:
		var request protocol.Request
		err := json.Unmarshal(msg.Payload(), &request)

		if err != nil {
			logger.Error("unmarshal payload error, ", string(msg.Payload()))
			return
		}

		go processRequest(request)
		break
	case subResponseTopic:
		var response protocol.Response
		err := json.Unmarshal(msg.Payload(), &response)
		if err != nil {
			logger.Error("unmarshal payload error, ", string(msg.Payload()))
		}
		break
	}
}

func processRequest(request protocol.Request) {
	switch request.Method {
	case "remote_control":
		var err error

		var cmd protocol.RemoteControlCommand
		err = json.Unmarshal([]byte(request.Params), &cmd)
		if err != nil {
			logger.Error("unmarshal remote control command error:", err.Error())
			return
		}

		defer func() {
			response := protocol.RemoteControlResponse{
				CtrlResult: 1,
				Message:    "",
				RespTime:   time.Now().Unix(),
				DeviceId:   cmd.DeviceId,
				CtrlTarget: cmd.CtrlTarget,
				CtrlAction: cmd.CtrlAction,
			}

			if err != nil {
				response.CtrlResult = 0
				response.Message = err.Error()
				logger.Error("control device error", err.Error())
			}

			data, _ := json.Marshal(response)
			resp := protocol.Response{
				MsgId: request.MsgId,
				Code:  protocol.CODE_RESPONSE,
				Data:  string(data),
			}

			responseChan <- resp
		}()

		err = ControlDevice(cmd.DeviceId, cmd.CtrlTarget, cmd.CtrlAction)
		if err != nil {
			logger.Error("control device error:", err.Error())
			return
		}

	case "realtime_state":
		var rtd protocol.DeviceRealtimeData
		err := json.Unmarshal([]byte(request.Params), &rtd)
		if err != nil {
			logger.Error("unmarshal realtime data error:", err.Error())
			return
		}

	}
}

//reportDeviceState report device state
func reportDeviceState(state protocol.DeviceRealtimeData) {
	data, _ := json.Marshal(state)
	request := protocol.Request{
		Version: "1.0",
		MsgId:   getSequenceNo(),
		Method:  "realtime_state",
		Params:  string(data),
	}
	reportChan <- request
}

func routeMqttMessage() {
	go func() {
		for {
			value := <-reportChan
			data, _ := json.Marshal(value)
			flag := mqttAdaptor.Publish(reportTopic, data)
			if !flag {
				logger.Error("publish report message failed", string(data))
			}
		}
	}()

	go func() {
		for {
			resp := <-responseChan
			data, _ := json.Marshal(resp)
			flag := mqttAdaptor.Publish(responseTopic, data)
			if !flag {
				logger.Error("publish response message failed", string(data))
			}
		}
	}()
}
