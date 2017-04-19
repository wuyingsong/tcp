package tcp

import (
	"errors"
	"log"
	"net"
	"os"
	"time"
)

const (
	//tcp conn max packet size
	defaultMaxPacketSize = 1024 << 10 //1MB

	READ_CHAN_SIZE  = 10
	WRITE_CHAN_SIZE = 10
)

var (
	logger *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "", log.Lshortfile)
}

type TcpServer struct {
	//TCP address to listen on
	tcpAddr string

	//the listener
	listener *net.TCPListener

	//callback is an interface
	//it's used to process the connect establish, close and data receive
	callback CallBack
	protocol Protocol

	//if srv is shutdown, close the channel used to inform all session to exit.
	exitChan chan struct{}

	maxPacketSize uint32        //single packet max bytes
	deadLine      time.Duration //the tcp connection read and write timeout
	bucket        *TcpConnBucket
}

func NewTcpServer(tcpAddr string, callback CallBack, protocol Protocol) *TcpServer {
	return &TcpServer{
		tcpAddr:  tcpAddr,
		callback: callback,
		protocol: protocol,

		bucket:        newTcpConnBucket(),
		exitChan:      make(chan struct{}),
		maxPacketSize: defaultMaxPacketSize,
	}
}

func (srv *TcpServer) ListenAndServe() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", srv.tcpAddr)
	if err != nil {
		return err
	}
	ln, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

func (srv *TcpServer) Serve(l *net.TCPListener) error {
	srv.listener = l
	defer func() {
		if r := recover(); r != nil {
			log.Println("Serve error", r)
		}
		srv.listener.Close()
	}()
	go func() {
		for {
			srv.removeClosedTcpConn()
			time.Sleep(time.Millisecond * 10)
		}
	}()

	var tempDelay time.Duration
	for {
		select {
		case <-srv.exitChan:
			return errors.New("TpcServer Closed.")
		default:
		}
		conn, err := srv.listener.AcceptTCP()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			log.Println("ln error:", err.Error())
			return err
		}
		tempDelay = 0
		tcpConn := srv.newTcpConn(conn, srv.callback, srv.protocol)
		srv.bucket.Put(tcpConn.RemoteAddr(), tcpConn)
		srv.callback.OnConnected(tcpConn)
	}
}

func (srv *TcpServer) newTcpConn(conn *net.TCPConn, callback CallBack, protocol Protocol) *TcpConn {
	if callback == nil {
		// if the handler is nil, use srv handler
		callback = srv.callback
	}
	return NewTcpConn(conn, callback, protocol)
}

func (srv *TcpServer) Connect(ip string, callback CallBack, protocol Protocol) (*TcpConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ip)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	tcpConn := srv.newTcpConn(conn, callback, protocol)
	return tcpConn, nil

}

func (srv *TcpServer) Close() {
	defer srv.listener.Close()
	for _, c := range srv.bucket.GetAll() {
		if !c.IsClosed() {
			c.Close()
		}
	}
}

func (srv *TcpServer) removeClosedTcpConn() {
	for {
		select {
		case <-srv.exitChan:
			return
		default:
			removeKey := make(map[string]struct{})
			for key, conn := range srv.bucket.GetAll() {
				if conn.IsClosed() {
					removeKey[key] = struct{}{}
				}
			}
			for key, _ := range removeKey {
				srv.bucket.Delete(key)
			}
			time.Sleep(time.Millisecond * 10)
		}
	}
}
func (srv *TcpServer) GetAllTcpConn() []*TcpConn {
	conns := []*TcpConn{}
	for _, conn := range srv.bucket.GetAll() {
		conns = append(conns, conn)
	}
	return conns
}

func (srv *TcpServer) GetTcpConn(key string) *TcpConn {
	return srv.bucket.Get(key)
}
