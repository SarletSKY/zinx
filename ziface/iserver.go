package ziface

type IServer interface {
	Start() // 启动服务器方法
	Stop()
	Serve()                                 // 开启业务服务方法
	AddRouter(msgId uint32, router IRouter) // 添加路由功能
	GetConnMgr() IConnManager               // 获取连接管理
	SetOnConnStart(func(IConnection))       // 设置该Server连接创建时的需要回调的Hook函数
	SetOnConnStop(func(IConnection))        // Hook设置用来为用户在连接关闭前/开启后执行的自定义函数
	CallOnStart(conn IConnection)           // 调用连接创建之后需要回调的业务方法
	CallOnStop(conn IConnection)            // 调用连接停止之后需要回调的业务方法
}
