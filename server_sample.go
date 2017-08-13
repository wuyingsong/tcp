package main

import (
	"fmt"
	"awesomeProject/tcp"
	"log"
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
	getBytes := p.Bytes()

	fmt.Println("get msg ", string(getBytes))

	dataLen := len(getBytes)
	packet := tcp.NewDefaultPacket(tcp.PacketType(dataLen), getBytes)
	err := conn.Send(packet)
	if err != nil {
		log.Fatal(err)
	}
}

func main()  {
	tcpCallback := TcpCallback{}
	protocol := tcp.DefaultProtocol{}
	server := tcp.NewTCPServer("localhost:30000", &tcpCallback, &protocol)
	server.ListenAndServe()
}