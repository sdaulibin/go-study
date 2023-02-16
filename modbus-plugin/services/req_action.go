package services

import (
	"bufio"
	"encoding/hex"
	"modbus-plugin/constants"
	"modbus-plugin/logs"
	server_map "modbus-plugin/map"
	"modbus-plugin/modbus"
)

func Process(conn_key string) {
	conn := server_map.TcpConnMap[conn_key]
	server_map.TcpConnSyncMap[conn_key].Lock()
	//frame := []byte{0x64, 0x0C, 0x00, 0x01, 0x34, 0x00, 0x00, 0x00, 0x07, 0x61, 0xB0, 0x20}
	//frame := []byte{0x64, 0x30, 0x00, 0x01, 0x33, 0x00, 0x00, 0x00, 0x0A, 0x01, 0x00, 0x01, 0x22, 0x00, 0xFE, 0xFE, 0xFE, 0xFE, 0x68, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0x68, 0x14, 0x12, 0x34, 0x37, 0x33, 0x37, 0x35, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x34, 0xC4, 0xAB, 0x89, 0x67, 0x45, 0x39, 0x16, 0xA7, 0x8C, 0x20}
	frame := modbus.TcpFrame{}
	frame = frame.InitSendFrame(constants.FUNCID_COLLECT)
	frame = frame.SetAddrFrame(frame, []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31}, 4, 25)
	send := frame.GenTcpFrame(frame)
	conn.Write(send)
	reader := bufio.NewReader(conn)
	var recvByte [512]byte
	n, err := reader.Read(recvByte[:]) // 读取数据
	tf := RecvToTcpFrame(recvByte[:n])
	logs.Logger.Infof("trans []byte to frame :", tf)
	hexStr := hex.EncodeToString(recvByte[:n])
	logs.Logger.Infof("trans []byte to hex :", hexStr)
	if err != nil {
		logs.Logger.Infof("网关设备(id:" + conn_key + ")已断开连接!")
	}
	server_map.TcpConnSyncMap[conn_key].Unlock()
}

func RecvToTcpFrame(recv []byte) modbus.TcpFrame {
	tcpFrame := modbus.TcpFrame{}
	tcpFrame.InitRecvToTcpFrame(recv)
	return tcpFrame
}
