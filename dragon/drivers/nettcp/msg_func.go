package nettcp

import (
	"fmt"

	"errors"
)

// DeserializeFunc unmarshals bytes to message
type DeserializeFunc func([]byte) (Message, error)

var (
	// msgFuncFactory stores message handler to every msg type.
	msgFuncFactory map[uint16]DeserializeFunc
)

func init() {
	msgFuncFactory = map[uint16]DeserializeFunc{}
}

// Register register a deserializer for a message type
func Register(msgType uint16, deserializer func([]byte) (Message, error)) {
	if _, ok := msgFuncFactory[msgType]; ok {
		panic(fmt.Sprintf("trying to register message %d twice", msgType))
	}

	msgFuncFactory[msgType] = deserializer
}

// GetDeserializer get deserializer for msg type.
func GetDeserializer(msgType uint16) (DeserializeFunc, error) {
	if f, ok := msgFuncFactory[msgType]; ok {
		return f, nil
	}

	return nil, errors.New(fmt.Sprintf("DeserializeFunc for MessageType %d does not exists.", msgType))
}
