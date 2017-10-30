package main

import (
	"log"
	"time"

	"github.com/wuyingsong/tcp"
)

func main() {
	srv := tcp.NewAsyncTCPServer("localhost:9001", &EchoCallback{}, &EchoProtocol{})
	srv.SetReadDeadline(time.Second * 10)
	log.Println("start listen...")
	log.Println(srv.ListenAndServe())
	select {}
}
