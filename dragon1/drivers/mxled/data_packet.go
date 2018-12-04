package mxled

import (
	"errors"

	"yonghui.cn/dragon/utils"
)

type EnumFontSize byte

const (
	Size_16_16 EnumFontSize = iota //font size is 16*16
	Size_32_32                     //font size is 32*32
)

type EnumColor byte

const (
	Red EnumColor = iota
	Green
	Yello
)

type PacketBase struct {
	Header    []byte
	Address   byte
	FrameType byte
	Xor       byte
}

type DisplayPacket struct {
	PacketBase

	Length     byte
	FontSize   byte
	Color      byte
	DisplayWay byte
	PositionX  []byte
	PositionY  []byte
	Data       []byte
}

func NewDisplayPacket(
	addr byte,
	size EnumFontSize,
	color EnumColor,
	way byte,
	posX,
	posY,
	data []byte) *DisplayPacket {

	packet := &DisplayPacket{
		PacketBase: PacketBase{
			Header:    []byte{0x55, 0x55},
			Address:   addr,
			FrameType: 0xF1,
		},
		Length:     byte(len(data)),
		FontSize:   byte(size),
		Color:      byte(color),
		DisplayWay: way,
		PositionX:  posX,
		PositionY:  posY,
		Data:       data,
	}

	return packet
}

func (pkt *DisplayPacket) Package() []byte {
	buffer := make([]byte, 0)
	buffer = append(buffer, pkt.Header...)
	buffer = append(buffer, pkt.Address)
	buffer = append(buffer, pkt.FrameType)
	buffer = append(buffer, pkt.Length)
	buffer = append(buffer, pkt.FontSize)
	buffer = append(buffer, pkt.Color)
	buffer = append(buffer, pkt.DisplayWay)
	buffer = append(buffer, pkt.PositionX...)
	buffer = append(buffer, pkt.PositionY...)
	buffer = append(buffer, pkt.Data...)

	xor, err := utils.Xor(buffer)
	if err != nil {
		return buffer
	}

	buffer = append(buffer, xor)
	return buffer
}

type ReplyPacket struct {
	PacketBase
	State byte
}

func ParseToReplyPacket(buffer []byte) (*ReplyPacket, error) {
	if len(buffer) != 6 {
		return nil, errors.New("Invalid input:length must be 6")
	}

	//check header,must be 0x55 0x55
	if buffer[0] != 0x55 || buffer[1] != 0x55 {
		return nil, errors.New("Invalid input:must start with 0x55 0x55")
	}

	// checksum must match
	xor, _ := utils.Xor(buffer[0:5])
	if xor != buffer[5] {
		return nil, errors.New("Invalid input,checksum error")
	}

	// build ReplyPacket
	packet := &ReplyPacket{
		PacketBase: PacketBase{
			Header:    buffer[0:2],
			Address:   buffer[2],
			FrameType: buffer[3],
		},
		State: buffer[4],
	}
	return packet, nil
}
