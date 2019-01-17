package msg

import "encoding/json"

type GreetMessage struct {
	Data string
}

func (msg GreetMessage) MessageType() uint16 {
	return 0x0001
}

func (msg GreetMessage) Serialize() ([]byte, error) {
	buffer, err := json.Marshal(msg)
	return buffer, err
}
