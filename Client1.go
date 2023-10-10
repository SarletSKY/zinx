package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Println("Client Test ... start")
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	for {
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(1, []byte("Zinx V0.6 Client1 Test Message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err:", err)
			return
		}
		// 先读取headData
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) // 将msg填满
		if err != nil {
			fmt.Println("read head error")
			break
		}
		msgHead, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}
			fmt.Println("==> Recv Msg:ID=", msg.Id, ",len = ", msg.GetDataLen(), ",data = ", string(msg.Data))
		}
		time.Sleep(1 * time.Second)
	}
}
