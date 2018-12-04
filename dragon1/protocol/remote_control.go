package protocol

// RemoteControlCommand represents command to control device
type RemoteControlCommand struct {
	CtrlTarget string `json:"ctrl_target"`
	DeviceId   string `json:"device_id"`
	CtrlAction string `json:"ctrl_action"`
	CtrlTime   int64  `json:"ctrl_time"`
}

// RemoteControlResponse represents response for RemoteControlCommand
type RemoteControlResponse struct {
	CtrlResult int    `json:"ctrl_result"`
	Message    string `json:"message"`
	RespTime   int64  `json:"resp_time"`
	DeviceId   string `json:"device_id"`
	CtrlTarget string `json:"ctrl_target"`
	CtrlAction string `json:"ctrl_action"`
}
