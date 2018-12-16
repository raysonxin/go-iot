package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var (
	ErrNilData   = errors.New("nil data")
	ErrNotEnough = errors.New("not enough data")
)

const (
	MessageTypeBytes    = 2
	MessageLengthBytes  = 4
	MessageDataMaxBytes = 1 << 23
)

// LengthTypeDataCodec represents a message structure with length,type and data.
type LengthTypeDataCodec struct{}

// Encode encode the message to bytes
func (codec LengthTypeDataCodec) Encode(msg Message) ([]byte, error) {
	data, err := msg.Serialize()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, msg.MessageType())
	binary.Write(buf, binary.LittleEndian, uint32(len(data)))
	buf.Write(data)

	return buf.Bytes(), nil
}

// Decode decode bytes to message
// the structure of message is Length-Type-Data,
// the below code will analyze bytes according to the structure.
func (codec LengthTypeDataCodec) Decode(datas []byte) (Message, error) {
	if datas == nil {
		return nil, ErrNilData
	}

	if MessageLengthBytes+MessageTypeBytes >= len(datas) {
		return nil, ErrNotEnough
	}

	lenBytes := datas[0:MessageLengthBytes]
	lenBuf := bytes.NewReader(lenBytes)
	var msgLen uint32
	if err := binary.Read(lenBuf, binary.LittleEndian, &msgLen); err != nil {
		return nil, err
	}

	typeBytes := datas[MessageLengthBytes : MessageLengthBytes+MessageTypeBytes]
	typeBuf := bytes.NewReader(typeBytes)
	var typeValue uint16
	if err := binary.Read(typeBuf, binary.LittleEndian, &typeValue); err != nil {
		return nil, err
	}

	otherBytes := datas[MessageLengthBytes+MessageTypeBytes:]
	if uint32(len(otherBytes)) < msgLen {
		return nil, ErrNotEnough
	}

	dataBytes := otherBytes[:msgLen]
	f, err := GetDeserializer(typeValue)
	if err != nil {
		panic(err)
	}

	msg, err := f(dataBytes)

	return msg, err
}
