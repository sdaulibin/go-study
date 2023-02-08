package modbus

import server_map "modbus-plugin/map"

func InitTCPGo(conn_key string) {
	for {
		conn := server_map.TcpConnMap[conn_key]
		server_map.TcpConnSyncMap[conn_key].Lock()
		conn.Write([]byte{0x64, 0x0C, 0x00, 0x01, 0x34, 0x00, 0x00, 0x00, 0x07, 0x61, 0xB0, 0x20})
		server_map.TcpConnSyncMap[conn_key].Unlock()
	}
}
