package main

import (
	"log"

	"github.com/wuyingsong/tcp"
)

type callback struct{}

func (c *callback) OnMessage(conn *tcp.TCPConn, p tcp.Packet) {
	log.Println("server receive:", string(p.Bytes()[1:]))
}

func (c *callback) OnConnected(conn *tcp.TCPConn) {
	log.Println("new conn:", conn.GetRemoteAddr().String())
}

func (c *callback) OnDisconnected(conn *tcp.TCPConn) {
	log.Printf("%s disconnected\n", conn.GetRemoteIPAddress())
}

func (c *callback) OnError(err error) {
	log.Println(err)
}
