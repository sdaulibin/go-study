package test

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"modbus-plugin/services"
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
	frame := services.RecvToTcpFrame([]byte("A123456789123456789"))
	fmt.Printf("frame.Start: %v,%s\n", frame.Start, strconv.FormatInt(int64(frame.Start), 16))
	fmt.Printf("frame.FuncId: %v\n", frame.FuncId)

	data := make([]byte, 0)
	data = binary.BigEndian.AppendUint16(data, frame.FuncId)
	data = binary.BigEndian.AppendUint16(data, frame.FuncId)
	test.Log(hex.EncodeToString(data))

	i := 1
	fmt.Println(i<<7 + i)
}
