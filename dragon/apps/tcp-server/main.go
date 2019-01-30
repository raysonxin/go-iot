package main

import (
	"encoding/json"
	"fmt"
	"github.com/raysonxin/go-iot/dragon/apps/msg"
	"github.com/raysonxin/go-iot/dragon/drivers/tcp"
	"github.com/raysonxin/go-iot/dragon/utils"
	"github.com/sirupsen/logrus"
	"net"
)

func main() {

	log := utils.NewLogger("logs/tcp_log")

	tcp.Register(0x0001, func(bytes []byte) (tcp.Message, error) {
		var greet msg.GreetMessage
		err := json.Unmarshal(bytes, &greet)
		if err != nil {
			log.Error("unmarshal greet message err", err.Error())
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
		log.Info("recv tcp message")

		switch message.MessageType() {
		case 0x0001:
			log.Info("content: ", message.(msg.GreetMessage).Data)
		}
	})

	onCloseOption := tcp.OnCloseOption(func(socket tcp.Socket) {
		conn := socket.(*tcp.ServerConn)
		log.WithFields(logrus.Fields{
			"method":   "OnClose",
			"clientId": conn.Name(),
		}).Error("tcp client close")
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
