package main

import "github.com/wuyingsong/tcp"

func main() {
	protocol := &tcp.DefaultProtocol{}
	protocol.SetMaxPacketSize(100)
	go startServer(&callback{}, protocol)
	go startClient(&callback{}, protocol)
	select {}
}
