package ziface

type IMsgHandle interface {
	DoMsgHandler(request IRequest)          // 非阻塞执行消息
	AddRouter(msgId uint32, router IRouter) // 为消息队列添加处理逻辑
	StartWorkerPool()                       // 开启Worker工作池
	SendMsgToTaskQueue(request IRequest)    // 将消息送进TaskQueue,再由Worker处理
}
