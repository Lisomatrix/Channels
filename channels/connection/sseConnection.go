package connection

import "github.com/gorilla/websocket"

// SSEConnection - Implementation of the Connection interface for Server Sent Event clients
type SSEConnection struct {
	messageSendChannel chan []byte
	ws                 *websocket.Conn
	onCloseCB          func()
	onMessage          func([]byte)
	isClosed           bool
}

// Init - Initialize connection channels
func (connection *SSEConnection) Init() {
	connection.isClosed = false
	connection.messageSendChannel = make(chan []byte, 5)
}

// Send - Enqueue message into channel
func (connection *SSEConnection) Send(payload []byte) {
	connection.messageSendChannel <- payload
}

// SendText - Enqueue message into channel
func (connection *SSEConnection) SendText(payload []byte) {
	connection.messageSendChannel <- payload
}

// SetOnMessage - In this type of connection the callback won't do anything
func (connection *SSEConnection) SetOnMessage(cb func([]byte)) {

}

// SetOnClose - Callback for when connection closes
func (connection *SSEConnection) SetOnClose(cb func()) {

}

// Close - Close the current connection
func (connection *SSEConnection) Close() {
	connection.isClosed = true
}
