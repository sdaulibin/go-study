package modbus

type TcpFrame struct {
	Start  uint16 //帧起始符 固定为 0x64
	Length []byte //帧总长度 例如整个数据包有120个字节 ， 此处应填 0x78 0x00
	Fixed  uint16 //保留字段  固定为1
	Serial []byte //帧序号 同一条命令下行时携带的序号，对此数据回复时应在同样的位置携带此标志，表示对之前某一条命令的回复
	FuncId uint16 //标识符 标识功能ID，在具体功能里定义
	Data   []byte //数据区 用户区数据
	Crc    []byte //CRC16校验 2个字节帧格式校验
	End    uint16 //帧结束符 0x20
}
