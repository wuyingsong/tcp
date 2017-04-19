package tcp

type CallBack interface {
	OnConnected(conn *TcpConn)

	OnMessage(conn *TcpConn, p Packet)

	OnDisconnected(conn *TcpConn)
}
