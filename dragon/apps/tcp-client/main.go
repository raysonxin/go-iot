package main

import (
	"fmt"
	msg2 "github.com/raysonxin/go-iot/dragon/apps/msg"
	"github.com/raysonxin/go-iot/dragon/drivers/tcp"
	"net"
	"time"
)

func main() {
	tcpaddr, err := net.ResolveTCPAddr("tcp4", "localhost:8989")

	if err != nil {
		fmt.Println("resolve tcp address err" + err.Error())
		return
	}

	//DialTCP建立一个TCP连接
	//net参数是"tcp4"、"tcp6"、"tcp"
	//laddr表示本机地址，一般设为nil
	//raddr表示远程地址
	tcpconn, err2 := net.DialTCP("tcp", nil, tcpaddr)
	if err2 != nil {
		fmt.Println("dial tcp server err" + err2.Error())
		return
	}

	codec := tcp.NewLengthTypeDataCodec()

	for i := 0; i < 10; i++ {
		msg := msg2.GreetMessage{
			Data: "HelloWorld",
		}

		buffer, _ := codec.Encode(msg)

		//fmt.Println("send content: "+string(buffer))

		//向tcpconn中写入数据
		_, err3 := tcpconn.Write(buffer)
		if err3 != nil {
			fmt.Println("Write message err" + err3.Error())
			return
		}

		time.Sleep(3 * time.Second)
	}

}
