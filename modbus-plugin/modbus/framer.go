package modbus

import (
	"log"
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
	result := make([]byte, 0)
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)
	for i := 0; i < t.NumField(); i++ {
		log.Println(t.Field(i).Name)
		log.Println(v.FieldByName(t.Field(i).Name))
		log.Println(v.Field(i))
		//bytes.Join(data, result)
	}
	return result
}
