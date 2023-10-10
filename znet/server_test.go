package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func ClientTest() {
	fmt.Println("Client Test ...start")
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err exit !")
		return
	}

	for {
		_, err := conn.Write([]byte("hello ZINX"))
		if err != nil {
			fmt.Println("write error err", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
			return
		}
		fmt.Printf("server call back: %s,cnt = %d\n", buf, cnt)
		time.Sleep(1 * time.Second)
	}
}

func TestServer_Serve(t *testing.T) {
	// 1. 服务器端测试
	s := NewServer()
	// 2.客户端测试
	go ClientTest()
	// 3.开启服务
	s.Serve()
}
