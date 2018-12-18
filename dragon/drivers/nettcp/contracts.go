package nettcp

// Socket network socket
type Socket interface {
	Write(Message) error
	Close()
}

type onConnectFunc func(Socket)

type onMessageFunc func(Message, Socket)

type onCloseFunc func(Socket)

type onErrorFunc func(Socket)
