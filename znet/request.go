package znet

import "zinx/ziface"

type Request struct {
	conn ziface.IConnection
	//data []byte
	msg ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 请求中获取信息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
