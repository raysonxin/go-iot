package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/raysonxin/go-iot/terminal/drivers"
)

func main() {
	b, err := ioutil.ReadFile("device-cfg.json")
	if err != nil {
		fmt.Println("read file error", err)
		return
	}

	fmt.Println(string(b))

	var cfgs []drivers.IotDeviceOption
	err = json.Unmarshal(b, &cfgs)
	if err != nil {
		fmt.Println("failed to unmarshal device config file")
		return
	}

	for _, v := range cfgs {
		if v.Protocol == "modbus-tcp" {
			recvChan := make(chan map[string]string, 0)
			tcpDriver, err := drivers.NewModbusTcpDriver(&v, recvChan)

			if err != nil {
				fmt.Println("create device driver error,", err)
				continue
			}
			go func() {
				for {
					datas := <-recvChan
					for k, v := range datas {
						fmt.Println(k, "=", v)
					}
				}
			}()

			go tcpDriver.Start()
		}
	}

	time.Sleep(90 * time.Second)
}
