package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"yonghui.cn/dragon/utils"
)

var logger *logrus.Entry

func main() {

	logger = utils.NewLogger("main entry", "logs/log")

	rootpath := utils.GetCurrentDirectory()

	devCfgPath := fmt.Sprintf("%s/config/vlong-cfg.json", rootpath)

	err := LoadDeviceConf(devCfgPath)
	if err != nil {
		fmt.Println("load vlongsoft device config file error, ", err.Error())
		panic(err)
		//return
	}

	for _, dev := range terminals {
		for _, swt := range dev.Switches {
			fmt.Println("start to open ", swt.Name, "...")
			err = ControlDevice(dev.DeviceId, swt.Name, "Open")
			if err != nil {
				fmt.Println("open ", swt.Name, " ", err.Error())
			}

			time.Sleep(10 * time.Second)
			fmt.Println("start to close ", swt.Name, "...")
			err = ControlDevice(dev.DeviceId, swt.Name, "Close")
			if err != nil {
				fmt.Println("close ", swt.Name, " ", err.Error())
			}

		}
	}

	fmt.Println("exit after 5 seconds")
}
