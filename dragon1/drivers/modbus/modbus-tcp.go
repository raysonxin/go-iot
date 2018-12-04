// this is a wrapped driver for modbus tcp device.
package modbus

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/goburrow/modbus"
	"yonghui.cn/dragon/drivers"
	"yonghui.cn/dragon/utils"
)

//ModbusTcpDriver modbus tcp driver
type ModbusTcpDriver struct {
	drivers.Driver

	tcpHandler *modbus.TCPClientHandler //modbus tcp client handler
}

// NewModbusTcpDriver create a modbus tcp device driver instance
func NewModbusTcpDriver(name, host string, port int, slaveId byte, timeout int) *ModbusTcpDriver {

	//create tcp handler for modbus
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", host, port))
	handler.Timeout = time.Duration(timeout) * time.Second
	handler.SlaveId = slaveId
	handler.Logger = log.New(os.Stdout, "ModbusTcp: ", log.LstdFlags)

	driver := &ModbusTcpDriver{
		Driver: drivers.Driver{
			Name: name,
		},
		tcpHandler: handler,
	}
	return driver
}

// Connect connect to device
func (driver *ModbusTcpDriver) Connect() {

}

// Collect the method to interact with device
func (driver *ModbusTcpDriver) ReadHoldingRegister(holders []ModbusHoldingAllocation, callback func(vals map[string]string)) {

	err := driver.tcpHandler.Connect()
	if err != nil {
		fmt.Println("Connect to modbus tcp server error")
		return
	}
	defer driver.tcpHandler.Close()

	client := modbus.NewClient(driver.tcpHandler)
	if len(holders) <= 0 {
		return
	}

	kvs := make(map[string]string, len(holders))

	for _, value := range holders {
		results, err := client.ReadHoldingRegisters(value.RegisterAddr, value.RegisterCount)
		if err != nil {
			fmt.Println("read err:", err)
		} else {
			num, _ := utils.BytesToInt(results)
			kvs[value.Name] = strconv.Itoa(num)
		}
	}

	callback(kvs)
}

func (driver *ModbusTcpDriver) WriteRegisters() error {
	err := driver.tcpHandler.Connect()
	if err != nil {
		return err
	}
	defer driver.tcpHandler.Close()

	client := modbus.NewClient(driver.tcpHandler)
	_, err = client.WriteMultipleRegisters(1, 2, []byte{0, 1, 2, 3})
	return err
}

// Disconnect stop
func (driver *ModbusTcpDriver) Disconnect() {
	driver.tcpHandler.Close()
}
