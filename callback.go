package tcp

type CallBack interface {
	OnConnected(conn *TCPConn)

	OnMessage(conn *TCPConn, p Packet)

	OnDisconnected(conn *TCPConn)
}
