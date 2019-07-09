package main

import (
	"log"
	"net"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func main() {
	serverinfo := "localhost:30000"
	conn, err := net.Dial("tcp", serverinfo)
	if err != nil {
		log.Println("conn err", err)
		return
	}


	var buf bytes.Buffer
	writehead := make([]byte, 4)
	data := "hello world!!"
	binary.BigEndian.PutUint32(writehead, uint32(len(data)))
	binary.Write(&buf, binary.BigEndian, writehead)
	buf.WriteString(data)
	conn.Write(buf.Bytes())

	readhead := make([]byte, 4)
	_, err = io.ReadFull(conn, readhead)
	if err != nil {
		log.Fatal(err)
	}

	packetLength := binary.BigEndian.Uint32(readhead)
	sendBuf := make([]byte, packetLength)
	_, err = io.ReadFull(conn, sendBuf)
	if err != nil {
		fmt.Println("err ", err)
	}

	fmt.Println(string(sendBuf))
	conn.Close()
}
