package main

import (
	"encoding/json"
	"fmt"
	"github.com/raysonxin/go-iot/dragon/apps/msg"
	"github.com/raysonxin/go-iot/dragon/drivers/tcp"
	"net"
)

func main() {

	tcp.Register(0x0001, func(bytes []byte) (tcp.Message, error) {
		var greet msg.GreetMessage
		err := json.Unmarshal(bytes, &greet)
		if err != nil {
			fmt.Println("unmarshal greet message err", err.Error())
			return nil, err
		}

		return greet, err
	})

	onBufferSizeOption := tcp.OnBufferSizeOption(256)

	onConnectOption := tcp.OnConnectOption(func(socket tcp.Socket) {
		sc := socket.(*tcp.ServerConn)
		fmt.Println("on connect" + sc.Name())
		return
	})

	onMessageOption := tcp.OnMessageOption(func(message tcp.Message, socket tcp.Socket) {
		fmt.Println("on message")
		switch message.MessageType() {
		case 0x0001:
			fmt.Println(" content: " + message.(msg.GreetMessage).Data)
		}
	})

	onCloseOption := tcp.OnCloseOption(func(socket tcp.Socket) {
		fmt.Println("on close")
	})

	setCodecFuncOptions := tcp.SetCodecFuncOption(func() tcp.MessageCodec {
		return tcp.NewLengthTypeDataCodec()
	})

	server := tcp.NewServer(onBufferSizeOption, onMessageOption, onConnectOption, onCloseOption, setCodecFuncOptions)

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 8989))
	if err != nil {
		fmt.Println("create listener err" + err.Error())
	}

	err = server.Start(l)
	if err != nil {
		fmt.Println("start server error" + err.Error())
	}
}
