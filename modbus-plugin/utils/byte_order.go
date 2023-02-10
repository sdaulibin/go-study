package utils

import (
	"encoding/binary"
)

func LittleOrder(sorce int16, len int) []byte {
	result := make([]byte, len)
	binary.LittleEndian.PutUint16(result, uint16(sorce))
	return result
}
