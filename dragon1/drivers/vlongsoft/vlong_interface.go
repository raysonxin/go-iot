package vlongsoft

import (
	"encoding/xml"
	"time"

	"github.com/hooklift/gowsdl/soap"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type TerminalCommand struct {
	XMLName xml.Name `xml:"urn:NewInterface TerminalCommand"`

	UserName string `xml:"UserName,omitempty"`

	PassWord string `xml:"PassWord,omitempty"`

	TerminalID int32 `xml:"TerminalID,omitempty"`

	ControllerID int32 `xml:"ControllerID,omitempty"`

	Command int32 `xml:"Command,omitempty"`
}

type GetVariableValue struct {
	XMLName xml.Name `xml:"urn:NewInterface GetVariableValue"`

	UserName string `xml:"UserName,omitempty"`

	PassWord string `xml:"PassWord,omitempty"`

	VariableType int32 `xml:"VariableType,omitempty"`

	TerminalID int32 `xml:"TerminalID,omitempty"`

	VariableID int32 `xml:"VariableID,omitempty"`
}

type Response struct {
	XMLName xml.Name `xml:"urn:NewInterface Response"`

	ResultCode int32 `xml:"ResultCode,omitempty"`

	Message string `xml:"Message,omitempty"`
}

type GetVariableResponse struct {
	XMLName xml.Name `xml:"urn:NewInterface GetVariableResponse"`

	TerminalID int32 `xml:"TerminalID,omitempty"`

	VariableID int32 `xml:"VariableID,omitempty"`

	VariableName string `xml:"VariableName,omitempty"`

	Value float32 `xml:"Value,omitempty"`

	AlarmMin float32 `xml:"AlarmMin,omitempty"`

	AlarmMax float32 `xml:"AlarmMax,omitempty"`

	DataTime string `xml:"DataTime,omitempty"`

	ResultCode int32 `xml:"ResultCode,omitempty"`

	Message string `xml:"Message,omitempty"`
}

type NewInterfacePortType interface {

	/* Service definition of function ns2__TerminalCommand */
	TerminalCommand(request *TerminalCommand) (*Response, error)

	/* Service definition of function ns2__GetVariableValue */
	GetVariableValue(request *GetVariableValue) (*GetVariableResponse, error)
}

type newInterfacePortType struct {
	client *soap.Client
}

func NewNewInterfacePortType(client *soap.Client) NewInterfacePortType {
	return &newInterfacePortType{
		client: client,
	}
}

func (service *newInterfacePortType) TerminalCommand(request *TerminalCommand) (*Response, error) {
	response := new(Response)
	err := service.client.Call("TerminalCommand", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *newInterfacePortType) GetVariableValue(request *GetVariableValue) (*GetVariableResponse, error) {
	response := new(GetVariableResponse)
	err := service.client.Call("GetVariableValue", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
