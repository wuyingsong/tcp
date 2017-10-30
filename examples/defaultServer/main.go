package main

import "github.com/wuyingsong/tcp"

func main() {
	protocol := &tcp.DefaultProtocol{}
	protocol.SetMaxPacketSize(100)
	startServer(&callback{}, protocol)
	startClient(&callback{}, protocol)
	select {}
}
