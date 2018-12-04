package vlongsoft

import "github.com/hooklift/gowsdl/soap"

type WsdlClient struct {
	url      string
	username string
	password string
}

// NewWsdlClient create a wsdl client with username,password,url
func NewWsdlClient(user, pwd, url string) *WsdlClient {
	return &WsdlClient{
		username: user,
		password: pwd,
		url:      url,
	}
}

// TerminalCommand send control command to terminal
func (w *WsdlClient) TerminalCommand(terminalId int, controllerId int, cmd int) (resp *Response, err error) {

	client := soap.NewClient(w.url, soap.WithHTTPHeaders(w.getHeaders()))
	svr := NewNewInterfacePortType(client)

	req := &TerminalCommand{
		UserName:     w.username,
		PassWord:     w.password,
		TerminalID:   int32(terminalId),
		ControllerID: int32(controllerId),
		Command:      int32(cmd),
	}

	resp, err = svr.TerminalCommand(req)

	return
}

// GetVariableValue get variaber current value
func (w *WsdlClient) GetVariableValue(varType int, terminalId int, varTypeId int) (resp *GetVariableResponse, err error) {

	client := soap.NewClient(w.url, soap.WithHTTPHeaders(w.getHeaders()))
	svr := NewNewInterfacePortType(client)

	req := &GetVariableValue{
		UserName:     w.username,
		PassWord:     w.password,
		VariableType: int32(varType),
		TerminalID:   int32(terminalId),
		VariableID:   int32(varTypeId),
	}

	resp, err = svr.GetVariableValue(req)

	return
}

func (w *WsdlClient) getHeaders() map[string]string {
	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	return header
}
