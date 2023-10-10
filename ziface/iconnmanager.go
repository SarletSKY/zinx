package ziface

type IConnManager interface {
	Add(conn IConnection)                   // 添加连接
	Remove(conn IConnection)                // 删除连接
	Len() int                               // 获取连接
	Get(connId uint32) (IConnection, error) // 通过Id获取连接
	ClearConn()                             // 删除并清除所有连接
}
