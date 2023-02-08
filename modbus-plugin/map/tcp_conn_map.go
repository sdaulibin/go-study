package server_map

import (
	"net"
	"sync"
)

var TcpConnMap = make(map[string]net.Conn)        // 网关tcp连接集合
var TcpConnSyncMap = make(map[string]*sync.Mutex) //网关tcp连接互斥锁（保证在处理反馈的时候不受干扰）
