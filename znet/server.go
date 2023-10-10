package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name      string // 服务器名称
	IPVersion string // tcp4
	IP        string
	Port      int
	//Router    ziface.IRouter // server注册的连接对应的处理业务
	msgHandler  ziface.IMsgHandle
	connMgr     ziface.IConnManager
	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server listenner at IP:%s,Port %d, is starting\n", s.IP, s.Port)
	fmt.Printf("[Zinx]Version %s,MaxConn:%d,MaxPacketSize:%d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
	go func() {
		// 0. 启动工作池
		s.msgHandler.StartWorkerPool()
		// 1. 获取一个tcp的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tpc addr err :", err)
			return

		}
		//2. 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start Zinx server", s.Name, "sncc listening...")

		// 生成自动的id的方法
		var cid uint32
		cid = 0
		// 3.启动Server网络连接业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			// 如果连接管理数量大于全局环境变量
			if s.connMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}
			// TODO:Server.Start() 设置服务器最大连接控制，如果超过最大连接，则关闭新的连接
			// TODO:Server Start() 处理该新连接请求的业务方法，此时handler和conn应该是绑定的
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++
			// 开启当前连接的处理业务
			go dealConn.Start()
			/*go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf) // 监听的数据读取到buf中
					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}
					// 回显
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err", err)
						continue
					}
				}
			}()*/
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[Stop] Zinx server ,name", s.Name)
	s.connMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
	//TODO Server.serve() 如果在启动服务的时候还要处理其他的事情，则可以在这里添加
	select {} // 阻塞
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router) // 给对象的router添加
	fmt.Println("Add Router Succ!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.connMgr
}

func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart...")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop...")
		s.OnConnStop(conn)
	}
}

// NewServer 初始化设置环境变量
func NewServer() ziface.IServer {
	//新建全局环境
	utils.GlobalObject.Reload() // 给全剧变量重新设置值
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		connMgr:    NewConnManager(),
	}
	return s
}

func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallBackToClient")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}
