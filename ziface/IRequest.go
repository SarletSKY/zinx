package ziface

// IRequest 将客户端请求的连接和请求封装在request
type IRequest interface {
	GetConnection() IConnection // 获取请求连接信息
	GetData() []byte            // 获取连接数据
	GetMsgId() uint32
}
