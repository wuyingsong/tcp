package tcp

import (
	"errors"
	"net"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	ErrConnClosing = errors.New("use of closed network connection")
	ErrBufferFull  = errors.New("the async send buffer is full")
)

type TCPConn struct {
	callback CallBack
	protocol Protocol

	conn      *net.TCPConn
	readChan  chan Packet
	writeChan chan Packet

	exitChan  chan struct{}
	closeOnce sync.Once
	exitFlag  int32
	err       error
}

func NewTCPConn(conn *net.TCPConn, callback CallBack, protocol Protocol) *TCPConn {
	c := &TCPConn{
		conn:      conn,
		callback:  callback,
		protocol:  protocol,
		readChan:  make(chan Packet, readChanSize),
		writeChan: make(chan Packet, writeChanSize),
		exitChan:  make(chan struct{}),
		exitFlag:  1,
	}
	c.Serve()
	return c
}

func (c *TCPConn) Serve() {
	defer func() {
		if r := recover(); r != nil {
			logger.Println("tcp conn(%v) Serve error, %v ", c.RemoteIP(), r)
		}
	}()
	c.callback.OnConnected(c)
	go c.readLoop()
	go c.writeLoop()
	go c.handleLoop()
}

func (c *TCPConn) readLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.exitChan:
			return
		default:
			p, err := c.protocol.ReadPacket(c.conn)
			if err != nil {
				return
			}
			c.readChan <- p
		}
	}
}

func (c *TCPConn) ReadPacket() (Packet, error) {
	if c.protocol == nil {
		return nil, errors.New("no protocol impl")
	}
	return c.protocol.ReadPacket(c.conn)
}

func (c *TCPConn) writeLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.exitChan:
			return
		case p := <-c.writeChan:
			if p == nil {
				continue
			}
			if err := c.protocol.WritePacket(c.conn, p); err != nil {
				return
			}
		}
	}
}

func (c *TCPConn) handleLoop() {
	defer func() {
		recover()
		c.Close()
	}()
	for {
		select {
		case <-c.exitChan:
			return
		case p := <-c.readChan:
			if p == nil {
				continue
			}
			c.callback.OnMessage(c, p)
		}
	}
}

func (c *TCPConn) Send(p Packet) error {
	if c.IsClosed() {
		return ErrConnClosing
	}
	select {
	case c.writeChan <- p:
		return nil
	default:
		return ErrBufferFull
	}
}

func (c *TCPConn) Close() {
	c.closeOnce.Do(func() {
		c.callback.OnDisconnected(c)
		atomic.StoreInt32(&c.exitFlag, 0)
		close(c.exitChan)
		close(c.readChan)
		close(c.writeChan)
		c.conn.Close()
	})
}

func (c *TCPConn) GetRawConn() *net.TCPConn {
	return c.conn
}

func (c *TCPConn) IsClosed() bool {
	return atomic.LoadInt32(&c.exitFlag) == 0
}

func (c *TCPConn) LocalAddr() string {
	return c.conn.LocalAddr().String()
}

func (c *TCPConn) LocalIP() string {
	return strings.Split(c.LocalAddr(), ":")[0]
}

func (c *TCPConn) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *TCPConn) RemoteIP() string {
	return strings.Split(c.RemoteAddr(), ":")[0]
}
