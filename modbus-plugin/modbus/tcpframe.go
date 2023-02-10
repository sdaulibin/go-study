package modbus

/*
*
字节 code 数据格式 描述
0 4 Bytes HEX TCP/IP 识别包头
4 11 Bytes ASCII 集中器地址
15 10 Bytes HEX 发送数据包头
25 3 Bytes HEX 控制码
28 6 Bytes BCD 采集器地址
34 6 Bytes BCD 表地址
40 2 Bytes HEX CRC 检验码
42 0X45 HEX 结束字符
*/
type TcpFrame struct {
	Start     []byte //0 4 Bytes HEX TCP/IP 识别包头, 0X7B 01 00 16,4 Bytes
	AsciiAddr []byte //4 11 Bytes ASCII 集中器地址
	SendHead  []byte //10 Bytes HEX 发送数据包头
	RecvHead  []byte //10 Bytes HEX 返回数据包头
	FuncId    []byte //3 Bytes HEX 控制码
	CollAddr  []byte //6 Bytes BCD 采集器地址
	MeterAddr []byte //6 Bytes BCD 表地址
	Data      []byte //数据区 用户区数据
	Crc       []byte //2 Bytes HEX CRC 检验码
	End       []byte //0X45 HEX 结束字符
}

type SendTcpFrame struct {
	TcpFrame
	SendHead uint16 //10 Bytes HEX 发送数据包头
}

type RecvTcpFrame struct {
	TcpFrame
	RecvHead uint16 //10 Bytes HEX 返回数据包头
}
