package main

import (
	"log"
	"net"
	"bytes"
	"encoding/binary"
)

func main() {
	serverinfo := "localhost:30000"
	conn, err := net.Dial("tcp", serverinfo)
	if err != nil {
		log.Println("conn err", err)
		return
	}


	var buf bytes.Buffer
	head := make([]byte, 4)
	data := "hello world!!"
	binary.BigEndian.PutUint32(head, uint32(len(data)))
	binary.Write(&buf, binary.BigEndian, head)
	buf.WriteString(data)
	conn.Write(buf.Bytes())
	conn.Close()
}
