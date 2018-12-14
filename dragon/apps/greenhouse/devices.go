package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"yonghui.cn/dragon/protocol"

	"yonghui.cn/dragon/drivers/vlongsoft"
)

var terminals []vlongsoft.Terminal

func loadJsonConfig(cfgPath string, v interface{}) error {
	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		//logger.Fatal("load config file ", cfgPath, " error, ", err.Error())
		return err
	}

	fmt.Println(string(b))
	err = json.Unmarshal(b, &v)
	if err != nil {
		//logger.Fatal("unmarshal config file ", cfgPath, " error,", err.Error())
		return err
	}

	return nil
}

//loadDeviceConf load controller config file
func LoadDeviceConf(cfgPath string) error {
	return loadJsonConfig(cfgPath, &terminals)
}

// func StartCycleMonitor(interval int) {
// 	if len(terminals) <= 0 {
// 		logger.Fatal("startCycleMonitor error,no device configuration.")
// 		return
// 	}

// 	task = cron.New()
// 	span := fmt.Sprintf("@every 10%ds", interval)
// 	task.AddFunc(span, func() {
// 		for _, dev := range terminals {
// 			go collectState(dev)
// 		}
// 	})
// 	task.Start()
// }

// collectState used to get vlongsoft controllers's state
func collectState(terminal vlongsoft.Terminal) {
	rd := protocol.DeviceRealtimeData{
		DeviceId:    terminal.DeviceId,
		CaptureTime: time.Now().Unix(),
		Datas:       make([]protocol.ItemValue, 0),
	}

	for _, item := range terminal.Switches {
		state, err := terminal.GetSwitchState(item.Name)
		if err != nil {
			//	logger.Error("collect device states error, [device_id=", terminal.DeviceId, ",switch_id="+item.Name, "],error:", err.Error())
			continue
		}
		swt := protocol.ItemValue{
			Param: item.Name,
			Value: strconv.Itoa(state),
			Order: 0,
		}
		rd.Datas = append(rd.Datas, swt)
	}

	rd.ReportTime = time.Now().Unix()

	//	reportDeviceState(rd)
}

func ControlDevice(devId, target, action string) error {
	t, err := findTerminal(devId)
	if err != nil {
		return err
	}

	switch action {
	case "Open":
		err = t.Open(target)
		break
	case "Close":
		err = t.Close(target)
		break
	case "Stop":
		err = t.Stop(target)
		break
	default:
		err = errors.New("not support command")
	}
	return err
}

func findTerminal(devId string) (*vlongsoft.Terminal, error) {
	var terminal vlongsoft.Terminal
	flag := false
	for _, dev := range terminals {
		if dev.DeviceId == devId {
			terminal = dev
			flag = true
			break
		}
	}

	if !flag {
		return nil, errors.New("not found a device which device id is" + devId)
	}
	return &terminal, nil
}
