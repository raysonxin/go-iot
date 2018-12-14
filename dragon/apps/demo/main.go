package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"yonghui.cn/dragon/protocol"
	"yonghui.cn/dragon/utils"

	"yonghui.cn/dragon/drivers/modbus"
	"yonghui.cn/dragon/drivers/mqtt"
)

var adaptor *mqtt.Adaptor

func main() {

	//initLogger()

	//logrus.Warn("test abcd.")

	// logger := utils.NewLogger("main entry", "logs/log")
	// logger.Error("test error log")
	// logger.Info("this is info log")
	// logger.Warn("this is warn log")
	// logger.Fatal("this is fatal log")

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	adaptor = mqtt.NewAdaptorWithAuth("tcp://10.0.90.55:1883", "dragon-test", "admin", "admin")
	err := adaptor.Connect()
	if err != nil {
		panic(err)
	}
	fmt.Println("connect to mqtt broker")
	go adaptor.Subscribe("gpgs/gateway/response/#", func(msg mqtt.Message) {
		fmt.Println(msg.Topic() + ":" + string(msg.Payload()))
	})

	sendCmd("light", "Open")
	time.Sleep(15 * time.Second)
	sendCmd("light", "Close")

	sendCmd("juanlian", "Open")
	time.Sleep(5 * time.Second)
	sendCmd("juanlian", "Stop")
	time.Sleep(5 * time.Second)
	sendCmd("juanlian", "Close")
	time.Sleep(5 * time.Second)
	sendCmd("juanlian", "Stop")
	time.Sleep(5 * time.Second)
	// ret = adaptor.Publish("test/we", []byte("abcd"))
	// if !ret {
	// 	fmt.Println("publish:false")
	// }

	<-sigc
	adaptor.Finalize()
}

func sendCmd(target, action string) {
	cmd := protocol.RemoteControlCommand{
		CtrlTarget: target,
		DeviceId:   "1234",
		CtrlAction: action,
		CtrlTime:   time.Now().Unix(),
	}

	data, _ := json.Marshal(cmd)

	request := protocol.Request{
		Version: "1.0",
		MsgId:   1,
		Method:  "remote_control",
		Params:  string(data),
	}

	b, _ := json.Marshal(request)

	ret := adaptor.Publish("gpgs/gateway-sub/request/gateway-dev", b)
	if !ret {
		fmt.Println("publish:false")
	}
}

func testModbusTcp() {
	b, err := ioutil.ReadFile("device-cfg.json")
	if err != nil {
		fmt.Println("read file error", err)
		return
	}

	fmt.Println(string(b))

	var cfgs []modbus.IotDeviceOption
	err = json.Unmarshal(b, &cfgs)
	if err != nil {
		fmt.Println("failed to unmarshal device config file")
		return
	}

	for _, v := range cfgs {
		if v.Protocol == "modbus-tcp" {

			var addr modbus.ModbusTcpAddress
			err = utils.MapToStruct(v.Address, &addr)

			//recvChan := make(chan map[string]string, 0)
			tcpDriver := modbus.NewModbusTcpDriver("test-modbus-tcp", addr.Host, addr.Port, addr.SlaveId, v.Timeout)

			var holders []modbus.ModbusHoldingAllocation
			err = utils.MapToStruct(v.Params, &holders)

			go func() {
				for {
					tcpDriver.ReadHoldingRegister(holders, func(vals map[string]string) {
						fmt.Println(tcpDriver.Name + " ")
						for k, v := range vals {
							fmt.Println(k + "=" + v)
						}
					})
					time.Sleep(3 * time.Second)
				}
			}()

			//	go tcpDriver.Start()
		}
	}
}
