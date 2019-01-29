package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// ServerConn represents a server connection
type ServerConn struct {
	connId  int64
	belong  *Server
	rawConn net.Conn
	codec   MessageCodec

	once     *sync.Once
	wg       *sync.WaitGroup
	buffer   []byte
	sendCh   chan []byte
	handleCh chan Message

	mu     sync.Mutex
	name   string
	heart  int64
	ctx    context.Context
	cancel context.CancelFunc
}

// NewServerConn create a server connection
func NewServerConn(id int64, s *Server, c net.Conn) *ServerConn {
	sc := &ServerConn{
		connId:   id,
		belong:   s,
		rawConn:  c,
		once:     &sync.Once{},
		wg:       &sync.WaitGroup{},
		sendCh:   make(chan []byte, s.opts.bufferSize),
		handleCh: make(chan Message, s.opts.bufferSize),
		heart:    time.Now().UnixNano(),
		buffer:   make([]byte, 0),
	}
	sc.ctx, sc.cancel = context.WithCancel(context.WithValue(s.ctx, serverCtx, s))
	sc.name = c.RemoteAddr().String()
	return sc
}

// SetName set server connection's name
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

// SetHeartbeat set heartbeat update time
func (sc *ServerConn) SetHeartbeat(heartbeat int64) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.heart = heartbeat
}

func (sc *ServerConn) SetCodec(codec MessageCodec) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.codec = codec
}

// Start start server connection
func (sc *ServerConn) Start() {
	onConnect := sc.belong.opts.onConnect
	if onConnect != nil {
		onConnect(sc)
	}

	sc.wg.Add(1)
	go sc.readLoop()

	sc.wg.Add(1)
	go sc.writeLoop()

	sc.wg.Add(1)
	go sc.handleLoop()
}

// Close close server connection
func (sc *ServerConn) Close() {
	sc.once.Do(func() {
		onClose := sc.belong.opts.onClose
		if onClose != nil {
			onClose(sc)
		}

		sc.belong.conns.Delete(sc.connId)

		if sock, ok := sc.rawConn.(*net.TCPConn); ok {
			sock.SetLinger(0)
		}
		sc.rawConn.Close()

		sc.mu.Lock()
		sc.cancel()
		sc.mu.Unlock()

		close(sc.sendCh)
		close(sc.handleCh)
	})
}

// Write write message to client socket
func (sc *ServerConn) Write(msg Message) error {
	datas, err := sc.codec.Encode(msg)
	if err != nil {
		return err
	}

	select {
	case sc.sendCh <- datas:
		err = nil
	default:
		err = errors.New("ErrWouldBlock")
	}
	return err
}

// readLoop the loop method to handle read
func (sc *ServerConn) readLoop() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("read loop panic:", p)
		}
		sc.wg.Done()
		fmt.Println("readLoop go-routine exited.")
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
				return
			}

			sc.SetHeartbeat(time.Now().UnixNano())

			buffer = append(sc.buffer, buffer[0:n]...)
			sc.buffer = sc.codec.Decode(buffer, sc.handleCh)
		}
	}
}

// writeLoop the loop method to handle write
func (sc *ServerConn) writeLoop() {
	var pkt []byte

	defer func() {
		if p := recover(); p != nil {
			fmt.Println("writeLoop panic err")
		}

		sc.wg.Done()
		fmt.Println("writeLoop go-routine exited.")
		sc.Close()
	}()

	for {
		select {
		case <-sc.ctx.Done():
			return
		case <-sc.belong.ctx.Done():
			return
		case pkt = <-sc.sendCh:
			if pkt != nil {
				if _, err := sc.rawConn.Write(pkt); err != nil {
					fmt.Println("error writing data " + err.Error())
					return
				}
			}
		}
	}
}

// handleLoop the loop method to handle message
func (sc *ServerConn) handleLoop() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("handle loop panic error", p)
		}

		sc.wg.Done()
		fmt.Println("handleLoop go-routine exited.")
		sc.Close()
	}()

	for {
		select {
		case <-sc.ctx.Done():
			return
		case <-sc.belong.ctx.Done():
			return
		case msg := <-sc.handleCh:

			onMessage := sc.belong.opts.onMessage
			if onMessage != nil {
				onMessage(msg, sc)
			}
		}
	}
}
