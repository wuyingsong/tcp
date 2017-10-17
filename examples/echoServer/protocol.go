package main

import (
	"bufio"
	"io"

	"github.com/wuyingsong/tcp"
)

type EchoProtocol struct {
	data []byte
}

func (ep *EchoProtocol) Bytes() []byte {
	return ep.data
}

func (ep *EchoProtocol) ReadPacket(reader io.Reader) (tcp.Packet, error) {
	rd := bufio.NewReader(reader)
	bytesData, err := rd.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return &EchoProtocol{data: bytesData}, nil
}
func (ep *EchoProtocol) WritePacket(writer io.Writer, msg tcp.Packet) error {
	_, err := writer.Write(msg.Bytes())
	return err
}
