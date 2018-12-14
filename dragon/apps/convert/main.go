package main

import (
	"encoding/json"
	"fmt"

	"yonghui.cn/dragon/utils"
)

func main() {
	str := `{"Name":"abcd","Data":{"Age":30,"Job":"IT","Other":{"Test":"test"}}}`
	var target BaseStruct
	err := json.Unmarshal([]byte(str), &target)
	if err != nil {
		panic(err)
	}

	var data DataStruct
	err = utils.MapToStruct(target.Data, &data)
	if err != nil {
		panic(err)
	}

	fmt.Println(data.Age, " ", data.Job, " ", data.Other)
}

type BaseStruct struct {
	Name string
	Data interface{}
}

type DataStruct struct {
	Age   int
	Job   string
	Other interface{}
}

type OtherStruct struct {
	Test string
}
