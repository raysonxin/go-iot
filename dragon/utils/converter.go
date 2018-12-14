package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
)

// ParseToStruct convert string verb to struct
func ParseBytesToStruct(buffer []byte, v interface{}) error {

	//buffer, err := ioutil.ReadAll(r)
	err := json.Unmarshal(buffer, v)
	return err
}

func BytesToInt(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp int8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp int16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp int32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}

func UInt16ToBytesBigEndian(v uint16) []byte {
	return []byte{byte(v >> 8), byte(v)}
}

func IntToBytesBigEndian(v int) []byte {
	return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
}

//MapToStruct parse map to struct
func MapToStruct(m interface{}, v interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}

// Xor xor byte by byte
func Xor(buffer []byte) (byte, error) {
	if buffer == nil || len(buffer) == 0 {
		return 0, errors.New("the buffer is empty.")
	}

	ret := buffer[0]
	for i := 1; i < len(buffer); i++ {
		ret ^= buffer[i]
	}
	return ret, nil
}

func BytesToUInt32LittleEndian(datas []byte, i uint32) error {
	buff := bytes.NewReader(datas)
	//var i uint32
	if err := binary.Read(buff, binary.LittleEndian, &i); err != nil {
		return err
	}
	return nil
}
