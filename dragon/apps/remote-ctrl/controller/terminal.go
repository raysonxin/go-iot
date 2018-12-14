package controller

import (
	"errors"
	"strconv"

	"yonghui.cn/dragon/apps/remote-ctrl/soapsvr"
)

const (
	SWITCH_ON  = 1
	SWITCH_OFF = 0
)

type Terminal struct {
	TerminalId int                 `json:"terminal_id"`
	ServiceUrl string              `json:"service_url"`
	Username   string              `json:"username"`
	Password   string              `json:"password"`
	Switches   []Switch            `json:"switches"`
	handle     *soapsvr.WsdlClient `json:"-"`
}

type Switch struct {
	Name        string       `json:"name"`
	State       int          // current state:1-close,2-open
	Type        string       `json:"type"` // switch type:-single,-double
	Controllers []Controller `json:"controllers"`
}

type Controller struct {
	Name         string `json:"name"`
	ControllerId int    `json:"controller_id"`
}

func (c *Terminal) findSwitch(name string) (swt Switch, err error) {
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

func (c *Terminal) Open(name string) (err error) {
	swt, err := c.findSwitch(name)
	if err != nil {
		return
	}

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

func (c *Terminal) Close(name string) (err error) {
	swt, err := c.findSwitch(name)
	if err != nil {
		return
	}
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

func (c *Terminal) Stop(name string) (err error) {
	swt, err := c.findSwitch(name)
	if err != nil {
		return
	}

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

func (c *Terminal) openSingle(controllerId int) (err error) {
	c.checkWsdlClient()

	resp, err := c.handle.TerminalCommand(c.TerminalId, controllerId, SWITCH_ON)
	if err != nil {
		return
	}

	if resp.ResultCode != 5000 {
		code := strconv.Itoa(int(resp.ResultCode))
		return errors.New("error code" + code)
	}

	return
}

func (c *Terminal) closeSingle(controllerId int) (err error) {

	c.checkWsdlClient()

	resp, err := c.handle.TerminalCommand(c.TerminalId, controllerId, SWITCH_OFF)
	if err != nil {
		return
	}

	if resp.ResultCode != 5000 {
		code := strconv.Itoa(int(resp.ResultCode))
		return errors.New("error code" + code)
	}

	return
}

func (c *Terminal) checkWsdlClient() {
	if c.handle == nil {
		handle := soapsvr.NewWsdlClient(c.Username, c.Password, c.ServiceUrl)
		c.handle = handle
	}
}
