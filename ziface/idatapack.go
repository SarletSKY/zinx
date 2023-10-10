package ziface

type IDataPack interface {
	GetHeadLen() uint32 // 获取包头长度
	Pack(msg IMessage) ([]byte, error)
	UnPack([]byte) (IMessage, error)
}
