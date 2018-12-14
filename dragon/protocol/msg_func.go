package protocol

import "fmt"

// DeserializeFunc unmarshals bytes to message
type DeserializeFunc func([]byte) (Message, error)

var (
	// msgFuncFactory stores message handler to every msg type.
	msgFuncFactory map[uint16]DeserializeFunc
)

func init() {
	msgFuncFactory = map[uint16]DeserializeFunc{}
}

func Register(msgType uint16, deserializer func([]byte) (Message, error)) {
	if _, ok := msgFuncFactory[msgType]; ok {
		panic(fmt.Sprintf("trying to register message %d twice", msgType))
	}

	msgFuncFactory[msgType] = deserializer
}
