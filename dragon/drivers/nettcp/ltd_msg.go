package nettcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	ErrNilData   = errors.New("nil data")
	ErrNotEnough = errors.New("not enough data")
	ErrInvalid   = errors.New("invalid data")
)

const (
	MessageTypeBytes    = 2
	MessageLengthBytes  = 4
	MessageDataMaxBytes = 1 << 23
)

// LengthTypeDataCodec represents a message structure with length,type and data.
type LengthTypeDataCodec struct {
	buffer []byte
}

func NewLengthTypeDataCodec() LengthTypeDataCodec {
	codec := LengthTypeDataCodec{
		buffer: make([]byte, 2*BufferSize1024),
	}
	return codec
}

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
func (codec LengthTypeDataCodec) Decode(datas []byte, msgs []Message) error {

	if datas == nil {
		return nil
	}

	if MessageLengthBytes+MessageTypeBytes >= len(datas) {
		codec.buffer = append(codec.buffer, datas...)
		return nil
	}

	lenBytes := datas[0:MessageLengthBytes]
	lenBuf := bytes.NewReader(lenBytes)
	var msgLen uint32
	if err := binary.Read(lenBuf, binary.LittleEndian, &msgLen); err != nil {

		return ErrInvalid
	}

	typeBytes := datas[MessageLengthBytes : MessageLengthBytes+MessageTypeBytes]
	typeBuf := bytes.NewReader(typeBytes)
	var typeValue uint16
	if err := binary.Read(typeBuf, binary.LittleEndian, &typeValue); err != nil {
		return nil, nil, err
	}

	otherBytes := datas[MessageLengthBytes+MessageTypeBytes:]
	if uint32(len(otherBytes)) < msgLen {
		return nil, datas, ErrNotEnough
	}

	if uint32(len(otherBytes)) == msgLen {

	}
	dataBytes := otherBytes[:msgLen]
	f, err := GetDeserializer(typeValue)
	if err != nil {
		fmt.Println("deserilizer for type " + string(typeValue) + " does not exists.")

		//panic(err)
		//return nil, err
	}

	msg, err := f(dataBytes)
	if err != nil {

	}

	return msg, err
}

//解包
func Unpack(buffer []byte, readerChannel chan Message) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+ConstHeaderLength+ConstSaveDataLength {
			break
		}
		if string(buffer[i:i+ConstHeaderLength]) == ConstHeader {
			messageLength := BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstSaveDataLength])
			if length < i+ConstHeaderLength+ConstSaveDataLength+messageLength {
				break
			}
			data := buffer[i+ConstHeaderLength+ConstSaveDataLength : i+ConstHeaderLength+ConstSaveDataLength+messageLength]
			readerChannel <- data

			i += ConstHeaderLength + ConstSaveDataLength + messageLength - 1
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}
