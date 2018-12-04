package vlongsoft

import (
	"errors"
	"fmt"
	"sync"
)

const (
	SWITCH_ACTION_ON  = 1
	SWITCH_ACTION_OFF = 0

	SWITCH_STATE_OPENING = 2
	SWITCH_STATE_CLOSING = 1
	SWITCH_STATE_STOPPED = 3
	SWITCH_STATE_FAULT   = 0

	VariableType_SENSOR = 1 // 1 is sensor type
	// VariableType_CONTROLLER 2 is controller type
	VariableType_CONTROLLER = 2 // 2 is controller type
)

// Terminal is a real device ,used to control some switch
type Terminal struct {
	DeviceId   string    `json:"device_id"`
	TerminalId int       `json:"terminal_id"`
	ServiceUrl string    `json:"service_url"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Switches   []*Switch `json:"switches"`
	handle     *WsdlClient
}

// Switch is a logic controller,used to control front-end sensor or device
type Switch struct {
	Name        string       `json:"name"`
	State       int          `json:"state"` // current state:1-close,2-open
	Type        string       `json:"type"`  // switch type:-single,-double
	Controllers []Controller `json:"controllers"`
	lock        sync.Mutex   //lock
}

type Controller struct {
	Name         string `json:"name"`
	ControllerId int    `json:"controller_id"`
}

func (c *Terminal) findSwitch(name string) (swt *Switch, err error) {
	flag := false
	for _, v := range c.Switches {
		if v.Name == name {
			swt = v
			flag = true
			break
		}
	}

	if !flag {
		err = errors.New("target not found")
	}

	return
}

// Open open switch with the specified name
func (c *Terminal) Open(name string) (err error) {
	swt, err := c.findSwitch(name)
	if err != nil {
		return
	}

	swt.lock.Lock()
	defer swt.lock.Unlock()

	if len(swt.Controllers) < 1 {
		err = errors.New("controller config error")
		return
	}

	if swt.Type == "single" {
		err = c.openSingle(swt.Controllers[0].ControllerId)
		return
	} else if swt.Type == "double" {

		open := Controller{}
		close := Controller{}

		for _, v := range swt.Controllers {
			if v.Name == "open" {
				open = v
			} else if v.Name == "close" {
				close = v
			}
		}

		err = c.closeSingle(close.ControllerId)
		if err != nil {
			return
		}

		err = c.openSingle(open.ControllerId)
		return
	}
	return
}

// Close close switch with the specified name
func (c *Terminal) Close(name string) (err error) {
	swt, err := c.findSwitch(name)
	if err != nil {
		return
	}

	swt.lock.Lock()
	defer swt.lock.Unlock()

	if swt.Type == "single" {

		err = c.closeSingle(swt.Controllers[0].ControllerId)
		return
	} else if swt.Type == "double" {

		open := Controller{}
		close := Controller{}

		for _, v := range swt.Controllers {
			if v.Name == "open" {
				open = v
			} else if v.Name == "close" {
				close = v
			}
		}

		err = c.closeSingle(open.ControllerId)
		if err != nil {
			return
		}

		err = c.openSingle(close.ControllerId)
		return
	}
	return
}

// Stop stop switch with the specified name
func (c *Terminal) Stop(name string) (err error) {
	swt, err := c.findSwitch(name)
	if err != nil {
		return
	}

	swt.lock.Lock()
	defer swt.lock.Unlock()

	if swt.Type != "double" {
		return
	}

	open := Controller{}
	close := Controller{}

	for _, v := range swt.Controllers {
		if v.Name == "open" {
			open = v
		} else if v.Name == "close" {
			close = v
		}
	}

	err = c.closeSingle(open.ControllerId)
	if err != nil {
		return
	}

	err = c.closeSingle(close.ControllerId)
	return
}

// GetSwitchState get switch state with the specified name
func (c *Terminal) GetSwitchState(name string) (state int, err error) {
	swt, err := c.findSwitch(name)
	if err != nil {
		return -1, err
	}

	swt.lock.Lock()
	defer swt.lock.Unlock()

	if swt.Type == "single" {
		err, state = c.getControllerValue(swt.Controllers[0].ControllerId)
		return
	}
	if swt.Type == "double" {
		open := Controller{}
		close := Controller{}

		for _, v := range swt.Controllers {
			if v.Name == "open" {
				open = v
			} else if v.Name == "close" {
				close = v
			}
		}

		err, openState := c.getControllerValue(open.ControllerId)
		if err != nil {
			return -1, err
		}

		err, closeState := c.getControllerValue(close.ControllerId)
		if err != nil {
			return -1, err
		}

		state = c.checkDoubleSwitchState(openState, closeState)
		return state, nil
	}
	return SWITCH_STATE_FAULT, nil
}

func (c *Terminal) checkDoubleSwitchState(openState, closeState int) (state int) {
	current := openState<<2 + closeState

	switch current {
	case 9:
		state = SWITCH_STATE_OPENING
	case 6:
		state = SWITCH_STATE_CLOSING
	case 5:
		state = SWITCH_STATE_STOPPED
	case 10:
		state = SWITCH_STATE_FAULT
	}
	return
}

// openSingle send command to terminal to open controller
func (c *Terminal) openSingle(controllerId int) (err error) {
	c.checkWsdlClient()

	resp, err := c.handle.TerminalCommand(c.TerminalId, controllerId, SWITCH_ACTION_ON)
	if err != nil {
		return
	}

	fmt.Println("DeviceID:", c.DeviceId, "TerminalID:", c.TerminalId, "ControllerID:", controllerId, "Command:ON", "return code:", resp.ResultCode)
	if resp.ResultCode != 5000 {
		return errors.New(fmt.Sprintf("error code:%d", resp.ResultCode))
	}
	return
}

// closeSingle send command to terminal to close one controller
func (c *Terminal) closeSingle(controllerId int) (err error) {

	c.checkWsdlClient()

	resp, err := c.handle.TerminalCommand(c.TerminalId, controllerId, SWITCH_ACTION_OFF)
	if err != nil {
		return
	}

	if resp.ResultCode != 5000 {
		//code := strconv.Itoa(int(resp.ResultCode))
		//return errors.New("error code:" + code)
		return errors.New(fmt.Sprintf("error code:%d", resp.ResultCode))
	}

	return
}

// getControllerValue get the specified controller value
func (c *Terminal) getControllerValue(controllerId int) (err error, value int) {
	// check proxy instance
	c.checkWsdlClient()

	//invoke proxy request to get current value
	resp, err := c.handle.GetVariableValue(VariableType_CONTROLLER, c.TerminalId, controllerId)
	if err != nil {
		return
	}

	// request success,5000 indicates success
	if resp.ResultCode == 5000 {
		return nil, int(resp.Value)
	} else {
		//code := strconv.Itoa(int(resp.ResultCode))
		//return errors.New("error code:" + code), -1
		//msg :=
		return errors.New(fmt.Sprintf("error code:%d", resp.ResultCode)), int(resp.Value)
	}
}

// checkWsdlClient check the proxy,if null,create.
func (c *Terminal) checkWsdlClient() {
	if c.handle == nil {
		handle := NewWsdlClient(c.Username, c.Password, c.ServiceUrl)
		c.handle = handle
	}
}
