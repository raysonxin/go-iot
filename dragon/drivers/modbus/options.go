package modbus

// Driver device driver base

type IotDeviceOption struct {
	Id       string      `json:"id"`       //device instance id
	Protocol string      `json:"protocol"` //protocol to interact with device
	Address  interface{} `json:"address"`  //param to connect device
	Interval int         `json:"interval"` //the recycle to collect
	Timeout  int         `json:"timeout"`  //interact wait time
	Params   interface{} `json:"params"`   //device capability
}

//ModbusTcpAddress address defination for modbus tcp
type ModbusTcpAddress struct {
	Host    string `json:"host"`     //device service host,it's an ip as usual
	Port    int    `json:"port"`     //device service port
	SlaveId byte   `json:"slave_id"` //device slave id
}

//ModbusRtuAddress address defination for modbus rtu
type ModbusRtuAddress struct {
	PortName string `json:"port_name"` //com port
	BaudRate int    `json:"baud_rate"`
	DataBits int    `json:"data_bits"`
	StopBits int    `json:"stop_bits"`
	Parity   int    `json:"parity"`
}

//ModbusHoldingAllocation device capability defination
type ModbusHoldingAllocation struct {
	Name          string `json:"name"`           //param name
	Code          int    `json:"code"`           //param code
	Index         int    `json:"index"`          //index number
	RegisterAddr  uint16 `json:"register_addr"`  //register start address
	RegisterCount uint16 `json:"register_count"` //register count
}
