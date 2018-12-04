package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"yonghui.cn/dragon/apps/remote-ctrl/controller"
)

func main() {

	b, err := ioutil.ReadFile("cfg.json")
	if err != nil {
		fmt.Println("read file error", err)
		return
	}

	fmt.Println(string(b))

	var t controller.Terminal
	err = json.Unmarshal(b, &t)
	if err != nil {
		fmt.Println("failed to unmarshal  config file")
		return
	}

	time.Sleep(10 * time.Second)

	err = t.Open("light")
	fmt.Println("open light ", err)
	time.Sleep(5 * time.Second)

	err = t.Close("light")
	fmt.Println("close light ", err)
	time.Sleep(3 * time.Second)

	err = t.Close("juanlian")
	fmt.Println("close juanlian ", err)
	time.Sleep(5 * time.Second)

	err = t.Stop("juanlian")
	fmt.Println("close juanlian ", err)
	time.Sleep(5 * time.Second)

	err = t.Open("juanlian")
	fmt.Println("open juanlian ", err)

	time.Sleep(10 * time.Second)
}
