package main

import (
	"bufio"
	"encoding/hex"
	"log"
	"modbus-plugin/initialize"
	"modbus-plugin/utils"
	"net"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

func main() {
	initialize.InitConfig()

	listen, err := net.Listen("tcp", viper.GetString("server.address"))
	if err != nil {
		log.Println("listen tcp server port failed, err: ", err, viper.GetString("server.address"))
	}
	log.Println("tcp server start success :", viper.GetString("server.address"))
	for {
		conn, err := listen.Accept() // 监听客户端的连接请求
		if err != nil {
			log.Println("accept conn failed, err: ", err)
			continue
		}
		procConn(conn)
	}
}
func procConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	var recvByte [128]byte
	sendByte := make([]byte, 0)
	n, err := reader.Read(recvByte[:])
	if err != nil {
		log.Println("read from conn failed, err: ", err)
		return
	}
	recvStr := string(recvByte[:n])
	log.Println("trans []byte to string :", recvStr)
	hexStr := hex.EncodeToString(recvByte[:n])
	log.Println("trans []byte to hex :", hexStr)

	firstByte := int64(recvByte[0])
	if strconv.FormatInt(firstByte, 16) == "64" {
		len := int(recvByte[1]) + int(recvByte[2])
		nextByte := recvByte[4:8]
		flag := int64(recvByte[8])
		log.Println("[]byte length : ", len)
		log.Println("[]byte nextByte : ", hex.EncodeToString(nextByte))
		log.Println("[]byte flag : ", strconv.FormatInt(flag, 16))
		if strconv.FormatInt(flag, 16) == "81" {
			sendByte = utils.BytesCombine(sendByte, []byte{0x64, 0x0C, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00})
			sendByte = utils.BytesCombine(sendByte, nextByte)
			crc16 := utils.ModbusCrc16(hex.EncodeToString(sendByte))
			log.Println("csc16 ====>>>>>:", crc16)
			crc16Byte, _ := hex.DecodeString(crc16)
			sendByte = utils.BytesCombine(sendByte, crc16Byte)
			sendByte = utils.BytesCombine(sendByte, []byte{0x20})
			log.Println("send message to device :", hex.EncodeToString(sendByte))
			conn.Write(sendByte)
		}
		time.Sleep(time.Second * time.Duration(5))
		conn.Write([]byte{0x64, 0x0C, 0x00, 0x01, 0x34, 0x00, 0x00, 0x00, 0x07, 0x61, 0xB0, 0x20})
	}

}
