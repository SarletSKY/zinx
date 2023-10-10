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
		msg, _ := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx V0.6 Client Test Message")))
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
	/*//1.创建一个封包对象 dp
	dp := znet.NewDataPack()

	//2.封装msg数据包
	msg1 := &znet.Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}
	msg2 := &znet.Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err:", err)
		return
	}
	// 将msg1与msg2粘一起
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)*/
	//select {}
	/*for {
		if _, err = conn.Write([]byte("hahaha")); err != nil {
			fmt.Println("write error err", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
			return
		}

		fmt.Printf("serer call back :%s,cnt=%d\n", buf, cnt)
		time.Sleep(time.Second * 1)
	}*/
}
