package main

import (
	"github.com/wuyingsong/tcp"
)

func startServer(callback tcp.CallBack, protocol tcp.Protocol) error {
	srv := tcp.NewAsyncTCPServer("localhost:9001", callback, protocol)
	return srv.ListenAndServe()
}
