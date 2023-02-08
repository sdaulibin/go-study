package main

import (
	"bufio"
	"encoding/hex"
	"log"
	"modbus-plugin/initialize"
	server_map "modbus-plugin/map"
	"modbus-plugin/services"
	"modbus-plugin/utils"
	"net"
	"strconv"
	"sync"
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
	register(conn)
	key := strconv.FormatInt(time.Now().Unix(), 10)
	server_map.TcpConnMap[key] = conn
	//delete(server_map.GwChannelMap, key)             // 删除原网关设备通道并新建通道
	//server_map.GwChannelMap[key] = make(chan int, 1) // 创建网关通道
	var s sync.Mutex
	server_map.TcpConnSyncMap[key] = &s
	services.Process(key)
}

// 处理注册请求
func register(conn net.Conn) {
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
			sendByte = utils.BytesCombine(sendByte, []byte{0x64, 0x0C, 0x00, 0x01})
			sendByte = utils.BytesCombine(sendByte, nextByte)
			sendByte = utils.BytesCombine(sendByte, []byte{0x01})
			crc16 := utils.ModbusCrc16(hex.EncodeToString(sendByte))
			log.Println("csc16 ====>>>>>:", crc16)
			crc16Byte, _ := hex.DecodeString(crc16)
			sendByte = utils.BytesCombine(sendByte, crc16Byte)
			sendByte = utils.BytesCombine(sendByte, []byte{0x20})
			log.Println("send message to device :", hex.EncodeToString(sendByte))
			conn.Write(sendByte)
		}
	}
}
