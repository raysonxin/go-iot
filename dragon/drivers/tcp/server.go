package tcp

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// Server define server struct
type Server struct {
	opts     options
	ctx      context.Context
	cancel   context.CancelFunc
	conns    *sync.Map
	wg       *sync.WaitGroup
	mu       sync.Mutex
	listener net.Listener
}

// NewServer create a new server
func NewServer(opt ...ServerOption) *Server {
	var opts options
	for _, o := range opt {
		o(&opts)
	}

	if opts.bufferSize <= 0 {
		opts.bufferSize = 256
	}

	s := &Server{
		opts:  opts,
		conns: &sync.Map{},
		wg:    &sync.WaitGroup{},
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	return s
}

// ConnsSize get current connection count
func (s *Server) ConnsSize() int {
	var sz int
	s.conns.Range(func(k, v interface{}) bool {
		sz++
		return true
	})
	return sz
}

// Start start tcp server with listner
func (s *Server) Start(l net.Listener) error {
	defer func() {
		l.Close()
	}()

	s.wg.Add(1)
	var tempDelay time.Duration

	for {
		rawConn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay >= max {
					tempDelay = max
				}
				select {
				case <-time.After(tempDelay):
				case <-s.ctx.Done():
				}
			}
			return err
		}

		tempDelay = 0

		sz := s.ConnsSize()
		if sz >= MaxConnections {
			//fmt.print too many conns
			rawConn.Close()
			continue
		}

		connId := time.Now().UnixNano()
		sc := NewServerConn(connId, s, rawConn)
		sc.SetName(sc.rawConn.RemoteAddr().String())

		if s.opts.setCodecFunc != nil {
			sc.SetCodec(s.opts.setCodecFunc())
		} else {
			sc.SetCodec(NewLengthTypeDataCodec())
		}

		s.conns.Store(connId, sc)

		s.wg.Add(1)
		go func() {
			sc.Start()
		}()

		fmt.Println("Accepted client ", sc.Name())

	}

	//	return nil
}

// Stop stop tcp server,release resource
func (s *Server) Stop() {
	s.mu.Lock()
	listener := s.listener
	s.listener = nil
	s.mu.Unlock()
	listener.Close()

	conns := map[int64]*ServerConn{}
	s.conns.Range(func(k, v interface{}) bool {
		i := k.(int64)
		c := v.(*ServerConn)
		conns[i] = c
		return true
	})
	s.conns = nil

	for _, c := range conns {
		c.rawConn.Close()
		fmt.Println("close client", c.Name())
	}

	s.mu.Lock()
	s.cancel()
	s.mu.Unlock()

	s.wg.Wait()
	fmt.Println("server stopped gracefully,bye.")
	os.Exit(0)
}
