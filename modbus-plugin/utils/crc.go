package utils

import (
	"encoding/hex"

	"github.com/deatil/go-crc16/crc16"
)

func ModbusCrc16(data string) string {
	crc16Hex, _ := hex.DecodeString(data)
	return crc16.ToHexString(crc16.ChecksumMODBUS(crc16Hex))
}
