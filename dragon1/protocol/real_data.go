package protocol

// DeviceRealtimeData represents device's real-time status
type DeviceRealtimeData struct {
	DeviceId    string      `json:"device_id"`
	CaptureTime int64       `json:"capture_time"`
	ReportTime  int64       `json:"report_time"`
	Datas       []ItemValue `json:"datas"`
}

// ItemValue represents device's item status
type ItemValue struct {
	Param string `json:"param"`
	Value string `json:"value"`
	Order int    `json:"order"`
}
