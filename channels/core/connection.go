package core

// Connection - Interface for connections
type Connection interface {
	//Init(ws *websocket.Conn)
	Send([]byte)
	SendText([]byte)
	SetOnMessage(func([]byte))
	SetOnClose(func())
	SetOnHeartBeat(func())
	Close()
	IsConnected() bool
}
