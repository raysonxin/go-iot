package protocol

//Request protocol request base struct
type Request struct {
	Version string      `json:"version"`
	MsgId   uint16      `json:"msg_id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

//Response protocol response base struct
type Response struct {
	MsgId uint16      `json:"msg_id"`
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
}
