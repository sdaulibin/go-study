package modbus

import (
	"bytes"
	"encoding/hex"
	"modbus-plugin/constants"
	"modbus-plugin/utils"
	"reflect"
)

type Framer interface {
	InitSendFrame([]byte) TcpFrame
	SetAddrFrame(asciiAddr, collAddr, meterAddr []byte)
	GenTcpFrame([]byte) []byte
	GenLittleAddr(int16) []byte
	InitRecvToTcpFrame([]byte) TcpFrame
}

func (tcpFrame *TcpFrame) InitSendFrame(funcId []byte) TcpFrame {
	sendFrame := TcpFrame{}
	sendFrame.Start = constants.FRAME_START
	sendFrame.SendHead = constants.SEND_HEAD
	sendFrame.FuncId = funcId
	return sendFrame
}

func (tcpFrame *TcpFrame) SetAddrFrame(sendFrame TcpFrame, addr []byte, coll, meter int16) TcpFrame {
	sendFrame.FuncId = constants.FUNCID_COLLECT
	sendFrame.AsciiAddr = addr
	collAddr := genLittleAddr(coll)
	sendFrame.CollAddr = collAddr
	meterAddr := genLittleAddr(meter)
	sendFrame.MeterAddr = meterAddr
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
	crc, _ := hex.DecodeString(utils.CcittCrc(hex.EncodeToString(bytes.Join(temps, []byte("")))))
	temps[len(temps)-2] = crc
	temps[len(temps)-1] = constants.FRRAME_END
	return bytes.Join(temps, []byte(""))
}

func genLittleAddr(source int16) []byte {
	addr, _ := hex.DecodeString(hex.EncodeToString(utils.LittleOrder(source, 6)))
	return addr
}

func (tcpFrame *TcpFrame) InitRecvToTcpFrame(recv []byte) TcpFrame {
	len := len(recv)
	recvFrame := TcpFrame{}
	recvFrame.Start = recv[0:4]
	recvFrame.AsciiAddr = recv[4:15]
	recvFrame.RecvHead = recv[15:25]
	recvFrame.FuncId = recv[25:28]
	// recvFrame.CollAddr = recv[28:34]
	// recvFrame.MeterAddr = recv[34:40]
	recvFrame.Crc = recv[len-3 : len-1]
	recvFrame.End = recv[len-1:]
	return recvFrame
}
