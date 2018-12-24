package nettcp

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type ServerConn struct {
	connId  int64
	belong  *Server
	rawConn net.Conn
	codec   MessageCodec

	once     *sync.Once
	wg       *sync.WaitGroup
	buffer   []byte
	sendCh   chan []byte
	handleCh chan DecodeResult

	mu     sync.Mutex
	name   string
	heart  int64
	ctx    context.Context
	cancel context.CancelFunc
}

func NewServerConn(id int64, s *Server, c net.Conn) *ServerConn {
	sc := &ServerConn{
		connId:   id,
		belong:   s,
		rawConn:  c,
		once:     &sync.Once{},
		wg:       &sync.WaitGroup{},
		sendCh:   make(chan []byte, s.opts.bufferSize),
		handleCh: make(chan DecodeResult, s.opts.bufferSize),
		heart:    time.Now().UnixNano(),
	}
	sc.ctx, sc.cancel = context.WithCancel(context.WithValue(s.ctx, serverCtx, s))
	sc.name = c.RemoteAddr().String()
	return sc
}

func (sc *ServerConn) SetName(name string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.name = name
}

func (sc *ServerConn) Name() string {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	name := sc.name
	return name
}

func (sc *ServerConn) SetCodec(codec MessageCodec) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.codec = codec
}

func (sc *ServerConn) Start() {
	onConnect := sc.belong.opts.onConnect
	if onConnect != nil {
		onConnect(sc)
	}
}

func (sc *ServerConn) Close() {

}

func (sc *ServerConn) Write(msg Message) error {
	return nil
}

func (sc *ServerConn) readLoop() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("panic:", p)
		}
		sc.wg.Done()
		sc.Close()
	}()

	for {
		select {
		case <-sc.ctx.Done():
			fmt.Println("conn cancel")
			return
		case <-sc.belong.ctx.Done():
			fmt.Println("server cancel")
			return
		default:
			buffer := make([]byte, BufferSize1024)
			n, err := sc.rawConn.Read(buffer)
			if err != nil {
				sc.Close()
			}

			buffer = append(sc.buffer, buffer[0:n]...)
			sc.buffer = sc.codec.Decode(buffer, sc.handleCh)
		}
	}
}

func writeLoop(c Socket, wg *sync.WaitGroup) {

}

func (sc *ServerConn) handleLoop() {

}
