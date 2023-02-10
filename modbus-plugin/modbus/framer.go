package modbus

import (
	"bytes"
	"encoding/hex"
	"modbus-plugin/constants"
	"reflect"
)

type Framer interface {
	InitFrame() TcpFrame
	GenTcpFrame([]byte) []byte
}

func (tcpFrame *TcpFrame) InitFrame() TcpFrame {
	sendFrame := TcpFrame{}
	sendFrame.Start = constants.FRAME_START
	sendFrame.End = constants.FRRAME_END
	return sendFrame
}

func (tcpFrame *TcpFrame) GenTcpFrame(input TcpFrame) []byte {
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)
	temps := make([][]byte, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		if v.Field(i).Type().String() == "[]uint8" {
			temp, _ := hex.DecodeString(hex.EncodeToString(v.Field(i).Bytes()))
			if len(temp) == 0 {
				continue
			}
			temps[i] = temp
		}
	}
	return bytes.Join(temps, []byte(""))
}
