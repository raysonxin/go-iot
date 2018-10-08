package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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
