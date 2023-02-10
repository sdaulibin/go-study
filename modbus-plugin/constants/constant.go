package constants

var FRAME_START []byte = []byte{0x7B, 0x01, 0x00, 0x16}

// 发送数据包头
var SEND_HEAD []byte = []byte{0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x53}

// 返回数据包头
var RECV_HEAD []byte = []byte{0x53, 0x53, 0x53, 0x53, 0x53, 0x53, 0x53, 0x53, 0x53, 0x42}
var FRRAME_END []byte = []byte{0x45}

// 控制码:采集请求
var FUNCID_COLLECT = []byte{0x4D, 0x4D, 0x4D}

// 控制码:采集请求-应答1
var FUNCID_RESP_1 = []byte{0x4D, 0x4D, 0x4D}

// 控制码:采集请求-应答2
var FUNCID_RESP_2 = []byte{0x44, 0x44, 0x44}