package mxled

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

var ticker *time.Ticker

type MxLedDriver struct {
	name    string
	host    string
	port    int
	address byte

	conn   *net.UDPConn
	locker sync.Mutex
}

func NewMxLedDriver(host string, port int, addr byte) *MxLedDriver {
	driver := &MxLedDriver{
		host:    host,
		port:    port,
		address: addr,
		name:    fmt.Sprintf("auto-%s:%d", host, port),
	}
	return driver
}

func (d *MxLedDriver) SetName(name string) {
	d.name = name
}

func (d *MxLedDriver) Name() string {
	return d.name
}

func (d *MxLedDriver) connect() (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", d.host, d.port))
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (d *MxLedDriver) Start() {

}

func (d *MxLedDriver) checkConn() {
	conn, err := d.connect()
	if err == nil {
		d.conn = conn
		return
	}

	ticker = time.NewTicker(10 * time.Second)
	//auto re-connect
	go func() {
		for t := range ticker.C {
			conn, err := d.connect()
			if err != nil {
				fmt.Println("ticker at ", t)
			} else {
				d.conn = conn
				ticker.Stop()
				break
			}
		}
	}()
}

func (d *MxLedDriver) Stop() {
	//ticker.Stop()
	//d.conn.Close()
}

func (d *MxLedDriver) DisplayContent(
	data []byte,
	size EnumFontSize,
	color EnumColor,
	posX,
	posY uint16) error {

	pkt := NewDisplayPacket(
		d.address,
		size,
		color,
		0x00, //stop
		[]byte{byte(posX >> 8), byte(posX)},
		[]byte{byte(posY >> 8), byte(posY)},
		data)

	buffer := pkt.Package()

	return d.send(buffer)
}

func (d *MxLedDriver) send(buffer []byte) error {

	d.locker.Lock()
	defer d.locker.Unlock()

	conn, err := d.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(buffer)
	if err != nil {
		return err
	}

	resp := make([]byte, 6)
	_, err = conn.Read(resp)
	if err != nil {
		return err
	}

	respPkt, err := ParseToReplyPacket(resp)
	if err != nil {
		return err
	}

	if respPkt.State != 0 {
		return errors.New("controller returns state:failed")
	}
	return nil
}
