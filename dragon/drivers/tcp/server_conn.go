package tcp

import (
	"context"
	"errors"
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
		buffer:   make([]byte, 0),
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

	sc.wg.Add(1)
	go sc.readLoop()

	sc.wg.Add(1)
	go sc.writeLoop()

	sc.wg.Add(1)
	go sc.handleLoop()
}

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

func (sc *ServerConn) readLoop() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("panic:", p)
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
			}

			buffer = append(sc.buffer, buffer[0:n]...)

			sc.buffer = sc.codec.Decode(buffer, sc.handleCh)
		}
	}
}

func (sc *ServerConn) writeLoop() {
	var pkt []byte

	defer func() {
		if p := recover(); p != nil {
			fmt.Println("panic err")
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

func (sc *ServerConn) handleLoop() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("panic err")
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
		case msgHandler := <-sc.handleCh:
			//msgHandler
			f, err := GetDeserializer(msgHandler.Type)
			if err != nil {
				continue
			}

			msg, err := f(msgHandler.Datas)
			if err != nil {
				continue
			}

			onMessage := sc.belong.opts.onMessage
			if onMessage != nil {
				onMessage(msg, sc)
			}
		}
	}
}
