package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// Message represents the message structure contracts.
type Message interface {

	// MessageType to get message type
	MessageType() uint16

	// Serialize serialize Message into bytes.
	Serialize() ([]byte, error)
}

// MessageCodec represents the codec contracts.
type MessageCodec interface {

	// Encode encode the message to bytes
	Encode(Message) ([]byte, error)

	// Decode decode bytes to message
	Decode([]byte) (Message, error)
}

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

	//dataBytes := otherBytes[:msgLen]

	return nil, nil
}

// type Protocol interface {
// 	ReadPacket(conn *net.TCPConn) (Packet, error)
// }

// type Packet interface {
// 	Serialize() []byte
// }

// // TcpPacket packet struct fort tcp transport
// type TcpPacket struct {
// 	Length  int         //packet length
// 	PktType byte        //packet type:request or reply
// 	Content interface{} //packet content part
// }

// func (pkt *TcpPacket) Serialize() []byte {

// 	datas, err := json.Marshal(pkt.Content)
// 	if err != nil {
// 		return nil
// 	}

// 	buffer := make([]byte, 0)
// 	buffer = append(buffer, utils.IntToBytesBigEndian(len(datas)+1)...)
// 	buffer = append(buffer, pkt.PktType)
// 	buffer = append(buffer, datas...)

// 	return buffer
// }
