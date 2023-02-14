package test

import (
	"encoding/hex"
	"fmt"
	"modbus-plugin/constants"
	"modbus-plugin/modbus"
	"modbus-plugin/utils"
	"net"
	"strconv"
	"testing"
)

func Test_Main(test *testing.T) {
	conn, err1 := net.Dial("tcp", "127.0.0.1:8004")
	if err1 != nil {
		test.Error("connect to tcp server error :", err1)
		return
	}
	defer conn.Close()
	sendMsg, _ := hex.DecodeString("7177657177313233")
	test.Log("send msg to tcp server :", string(sendMsg))
	n, err2 := conn.Write(sendMsg) // 发送数据
	if err2 != nil {
		test.Error("send meg to tcp server error :", err2)
		return
	}
	test.Log("send msg to tcp server success :", n)
}

func Test_Uint16(test *testing.T) {
	test.Log(uint16([]byte("abc")[1]))
	test.Log(hex.EncodeToString([]byte("abc")))
	test.Log([]byte("abc")[1])
	test.Log(strconv.FormatInt(int64([]byte("abc")[1]), 16))
}

func Test_Crc(test *testing.T) {
	//data, _ := hex.DecodeString("313233")
	test.Log(utils.ModbusCrc16("640c00010000000016000000"))
}

func Test_Frame(test *testing.T) {
	tcpFrame := modbus.TcpFrame{}
	tcpFrame = tcpFrame.InitSendFrame(constants.FUNCID_COLLECT)
	tcpFrame = tcpFrame.SetAddrFrame(tcpFrame, []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31}, 4, 25)
	b := tcpFrame.GenTcpFrame(tcpFrame)
	fmt.Println(hex.EncodeToString(b))
	i := 1
	fmt.Println(i<<7 + i)
}
