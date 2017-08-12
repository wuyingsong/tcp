package main

import (
	"fmt"
	"awesomeProject/tcp"
)

type TcpCallback struct {

}

func (tcpCallbak *TcpCallback)OnConnected(conn *tcp.TCPConn) {
	fmt.Println("connected ",conn.LocalIP())

}
func (tcpCallbak *TcpCallback)OnDisconnected(conn *tcp.TCPConn) {
	fmt.Println("disconnected ",conn.LocalIP())
}

func (tcpCallbak *TcpCallback)OnMessage(conn *tcp.TCPConn, p tcp.Packet){
	fmt.Println("get msg ", string(p.Bytes()))
}

func main()  {
	tcpCallback := TcpCallback{}
	protocol := tcp.DefaultProtocol{}
	server := tcp.NewTCPServer("localhost:30000", &tcpCallback, &protocol)
	server.ListenAndServe()
}