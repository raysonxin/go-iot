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

// definitions about some constants.
const (
	MaxConnections    = 1000
	BufferSize128     = 128
	BufferSize256     = 256
	BufferSize512     = 512
	BufferSize1024    = 1024
	defaultWorkersNum = 20
)

// ContextKey is the key type for putting context-related data.
type contextKey string

// Context keys for messge, server and net ID.
const (
	messageCtx contextKey = "message"
	serverCtx  contextKey = "server"
	netIDCtx   contextKey = "netid"
)
