package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	TcpServer  ziface.IServer // 当前conn属于哪个Server，在conn初始化的时候添加即可
	Conn       *net.TCPConn
	ConnID     uint32
	isClosed   bool
	MsgHandler ziface.IMsgHandle
	//Router   ziface.IRouter
	//handleAPI    ziface.HandFunc
	ExitBuffChan chan bool
	msgChan      chan []byte // 无缓冲通道，同于读写之间的消息通道
	msgBuffChan  chan []byte // 有缓冲通道，同于读写之间的消息通道
	// 连接属性
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		//handleAPI:    callback_api,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:     make(map[string]interface{}),
	}
	// 为服务中的连接管理添加连接
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn reader exit!]")
	defer c.Stop()

	for {
		/*buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			c.ExitBuffChan <- true
			continue
		}*/

		// 创建数据包
		dp := NewDataPack()

		// 读取客户端的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			c.ExitBuffChan <- true
			break
		}
		// 拆包
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			c.ExitBuffChan <- true
			continue
		}
		var data []byte
		// 第二次再次获取
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitBuffChan <- true
				break
			}
		}
		msg.SetData(data)
		// 获取用户的request数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 通过工作池的队列去执行
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 通过信息句柄去执行
			go c.MsgHandler.DoMsgHandler(&req)
		}
		// 通过协程制定request
		/*go func(request ziface.IRequest) {
			//c.Router.PreHandle(request)
			c.Router.Handle(request)
			//c.Router.PostHandle(request)
		}(&req)*/
		/*		// 调用当前连接业务(这里执行的是当前conn绑定的handle方法)
				if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
					fmt.Println("connID", c.ConnID, "handle is error")
					c.ExitBuffChan <- true
					return
				}*/
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Write Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:,", err, "Conn Write exit")
				return
			}
		case data, ok := <-c.msgBuffChan:
			// 针对有缓冲channel需要进行数据处理
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:,", err, "conn Write exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()
	c.TcpServer.CallOnStart(c) // 调用Hook函数，在连接之后
	/*for {
		select {
		case <-c.ExitBuffChan:

			return //得到exit消息，不再阻塞
		}
	}*/
}

func (c *Connection) Stop() {
	if c.isClosed == true { // 已经关闭
		return
	}
	c.isClosed = true

	c.TcpServer.CallOnStop(c)
	c.Conn.Close() // 关闭socket
	c.ExitBuffChan <- true
	c.TcpServer.GetConnMgr().Remove(c) // 先关闭连接管理的连接
	close(c.ExitBuffChan)              // 关闭所有的管道
	close(c.msgBuffChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New(" Connect closed when send msg ")
	}

	// 将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack err msgId = ", msgId)
		return errors.New(" Pack error msg ")
	}
	/*// 写回客户端
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("Write msg id", msgId, "error")
		c.ExitBuffChan <- true
		return errors.New(" conn Write error")
	}*/
	c.msgChan <- msg
	return nil
}

func (c *Connection) SendBufMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg ")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msgId = ", msgId)
		return errors.New(" Pack error msg ")
	}
	c.msgChan <- msg
	return nil
}
