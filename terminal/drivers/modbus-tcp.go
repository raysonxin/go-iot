// this is a wrapped driver for modbus tcp device.
package drivers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/goburrow/modbus"
	"github.com/raysonxin/go-iot/utils"
	"github.com/robfig/cron"
)

//ModbusTcpDriver modbus tcp driver
type ModbusTcpDriver struct {
	Driver
	Options    *IotDeviceOption         //device options
	TcpHandler *modbus.TCPClientHandler //modbus tcp client handler

	task               *cron.Cron                //timer
	outputChannel      chan<- map[string]string  //output channel
	registerAllocation []ModbusHoldingAllocation //register defination
	holder             map[string]string         //temp value-container
}

// NewModbusTcpDriver create a modbus tcp device driver instance
func NewModbusTcpDriver(opts *IotDeviceOption, ch chan<- map[string]string) (*ModbusTcpDriver, error) {
	//parse modbus tcp address
	var addr ModbusTcpAddress
	err := utils.MapToStruct(opts.Address, &addr)
	if err != nil {
		return nil, errors.New("Modbus tcp address parse error.")
	}

	//create modbus handler then set param
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", addr.Host, addr.Port))
	handler.Timeout = time.Duration(opts.Timeout) * time.Second
	handler.SlaveId = addr.SlaveId
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	err = handler.Connect()
	if err != nil {
		return nil, err //errors.New("Connect to modbus tcp server error")
	}

	//parse modbus register allocation
	var holdAllocation []ModbusHoldingAllocation
	err = utils.MapToStruct(opts.Params, &holdAllocation)

	if err != nil {
		return nil, errors.New("Modbus tcp params parse error.")
	}

	driver := &ModbusTcpDriver{
		Driver: Driver{
			InstanceId: opts.Id,
		},
		TcpHandler:         handler,
		task:               cron.New(),
		Options:            opts,
		registerAllocation: holdAllocation,
		outputChannel:      ch,
	}
	return driver, nil
}

// Start start a timer to collect device data
func (driver *ModbusTcpDriver) Start() {
	driver.holder = make(map[string]string)

	//build spec
	spec := fmt.Sprintf("@every %ds", driver.Options.Interval)
	driver.task.AddFunc(spec, driver.collect)
	driver.task.Start()
}

// collect the method to interact with device
func (driver *ModbusTcpDriver) collect() {

	client := modbus.NewClient(driver.TcpHandler)
	if len(driver.registerAllocation) <= 0 {
		return
	}

	for _, value := range driver.registerAllocation {
		fmt.Println("name=", value.Name, "code=", value.Code, "addr=", value.RegisterAddr, "count=", value.RegisterCount)
		results, err := client.ReadHoldingRegisters(value.RegisterAddr, value.RegisterCount)
		if err != nil {
			fmt.Println("read err:", err)
		} else {
			num, _ := utils.BytesToInt(results)
			driver.holder[value.Name] = strconv.Itoa(num)
		}
	}
	driver.outputChannel <- driver.holder
}

// Stop stop timer
func (driver *ModbusTcpDriver) Stop() {
	driver.task.Stop()
	driver.TcpHandler.Close()
}
