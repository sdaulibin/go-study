package services

import (
	"bufio"
	"encoding/hex"
	"log"
	server_map "modbus-plugin/map"
	"modbus-plugin/utils"
	"strconv"
	"time"
)

func Process(conn_key string) {
	conn := server_map.TcpConnMap[conn_key]
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
			server_map.TcpConnSyncMap[conn_key].Lock()
			conn.Write(sendByte)
		}
		time.Sleep(time.Second * time.Duration(2))
		conn.Write([]byte{0x64, 0x0C, 0x00, 0x01, 0x34, 0x00, 0x00, 0x00, 0x07, 0x61, 0xB0, 0x20})
		server_map.TcpConnSyncMap[conn_key].Unlock()
	}

}
