package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter // 基础路由
}

//func (this *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ...\n"))
//	if err != nil {
//		fmt.Println("call back ping ping ping error")
//	}
//}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	fmt.Println("recv from cient : msgId =", request.GetMsgId(), ",data =", string(request.GetData()))
	/*_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}*/
	err := request.GetConnection().SendMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter // 基础路由
}

func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxROuter Handle")
	fmt.Println("recv from cient : msgId =", request.GetMsgId(), ",data =", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("Hello Zinx Router V0.8"))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin is Called...")
	conn.SetProperty("Name", "Aceld")
	conn.SetProperty("Home", "https://github.com/aceld/zinx")
	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionLast(conn ziface.IConnection) {
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}
	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home = ", home)
	}
	fmt.Println("DoConnectionLast is Called...")
}

//func (this *PingRouter) PostHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping ...\n"))
//	if err != nil {
//		fmt.Println("call back ping ping ping error")
//	}
//}

func main() {
	/*listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	// 从客户端读取粘包的数据，进行解析
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server accept err:", err)
		}

		// 处理客户端的请求
		go func(conn net.Conn) {
			dp := znet.NewDataPack()
			for {
				// 读取数据
				heapData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, heapData) // 必须把所有数据填满读出来
				if err != nil {
					fmt.Println("read head error")
					break
				}
				// 解包，将数据流入msg中
				msgHead, err := dp.UnPack(heapData)
				if err != nil {
					fmt.Println("server unpack err:", err)
					return
				}

				// 读取字节流
				if msgHead.GetDataLen() > 0 {
					msg := msgHead.(*znet.Message)
					msg.Data = make([]byte, msg.GetDataLen())
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack data err:", err)
						return
					}
					fmt.Println("==> Recv Msg:Id=", msg.Id, ",len=", msg.DataLen, ",data=", string(msg.Data))
				}
			}
		}(conn)
	}*/
	s := znet.NewServer()
	s.SetOnConnStart(DoConnectionBegin) // 设置执行哪个hook函数
	s.SetOnConnStop(DoConnectionLast)
	s.AddRouter(0, &PingRouter{})      // 父类实现了接口
	s.AddRouter(1, &HelloZinxRouter{}) // 父类实现了接口
	s.Serve()                          // 启动服务，主线程阻塞，携程开启服务
}
