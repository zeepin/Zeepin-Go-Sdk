package common

import (
	"encoding/binary"
	"bytes"
)

func Reverse(v []byte) []byte{
	result := make([]byte, len(v))
	for i := 0; i < len(v); i++{
		result[i] = v[len(v) - i -1]
	}
	return result
}

func IntToByte(n int) uint8{

	tmp := int64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, &tmp)
	len := len(bytesBuffer.Bytes())
	return  bytesBuffer.Bytes()[len-1]
}
