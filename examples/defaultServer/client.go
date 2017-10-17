package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/wuyingsong/tcp"
)

func startClient(callback tcp.CallBack, protocol tcp.Protocol) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:9001")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	tc := tcp.NewTCPConn(conn, callback, protocol)
	log.Println(tc.Serve())
	i := 0
	for {
		if tc.IsClosed() {
			break
		}
		msg := fmt.Sprintf("hello %d", i)
		log.Println("client send: ", msg)
		tc.AsyncWritePacket(&tcp.DefaultPacket{Type: 1, Body: []byte(msg)})
		i++
		time.Sleep(time.Second)
	}
}
