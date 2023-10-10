package ziface

import (
	"net"
)

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn // 从当前连接获取原始的socket TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
	SendMsg(msgId uint32, data []byte) error
	SendBufMsg(msgId uint32, data []byte) error // 有缓冲
	SetProperty(key string, value interface{})  // 设置连接属性
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
}
type HandFunc func(*net.TCPConn, []byte, int) error
