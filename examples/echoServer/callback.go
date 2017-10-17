package main

import (
	"log"

	"github.com/wuyingsong/tcp"
)

type EchoCallback struct{}

func (ec *EchoCallback) OnConnected(conn *tcp.TCPConn) {
	log.Println("new conn: ", conn.GetRemoteIPAddress())
}
func (ec *EchoCallback) OnMessage(conn *tcp.TCPConn, p tcp.Packet) {
	log.Printf("receive: %s", string(p.Bytes()))
	conn.AsyncWritePacket(p)
}
func (ec *EchoCallback) OnDisconnected(conn *tcp.TCPConn) {
	log.Printf("%s disconnected\n", conn.GetRemoteIPAddress())
}

func (ec *EchoCallback) OnError(err error) {
	log.Println(err)
}
