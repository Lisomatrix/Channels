// This package holds the connection interface implementations to be used by project
package connection

import (
	"log"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// OWebSocketConnection - O stands for optimized
// Implementation of the Connection interface for WebSocket clients
type OWebSocketConnection struct {
	messageSendChannel chan []byte
	ws                 net.Conn
	onCloseCB          func()
	onMessage          func([]byte)
	onHb			   func()
	isClosed           bool
	pongTimer		   *time.Timer
	didReceivePong 	   bool
}

// IsConnected - If connection is still alive
func (connection *OWebSocketConnection) IsConnected() bool {
	return !connection.isClosed
}

// Init - Initialize connection and start reading and writing messages
func (connection *OWebSocketConnection) Init(ws net.Conn) {

	connection.ws = ws
	connection.messageSendChannel = make(chan []byte, 10)

	go connection.readMessages()
	go connection.writeMessages()
}

// Send - Send the data to the client
func (connection *OWebSocketConnection) Send(data []byte) {
	connection.messageSendChannel <- data
}

// SendText - Send the data to the client
func (connection *OWebSocketConnection) SendText(data []byte) {
	connection.messageSendChannel <- data
}

// SetOnMessage - Set user sent message callback
func (connection *OWebSocketConnection) SetOnMessage(cb func([]byte)) {
	connection.onMessage = cb
}

func (connection *OWebSocketConnection) SetOnHeartBeat(cb func()) {
	connection.onHb = cb
}

// SetOnClose - Set on closed connection handler
func (connection *OWebSocketConnection) SetOnClose(cb func()) {
	connection.onCloseCB = cb
}

// Close - Close the current connection
func (connection *OWebSocketConnection) Close() {

	_ = connection.ws.Close()

	connection.isClosed = true

	if connection.onCloseCB != nil {
		connection.onCloseCB()
	}
}

func (connection *OWebSocketConnection) readMessages() {
	defer func() {
		connection.Close()
	}()

	msg := make([]wsutil.Message, 0, 4)

	for {

		if connection.isClosed {
			return
		}

		msg, err := wsutil.ReadClientMessage(connection.ws, msg[:0])
		if err != nil {
			log.Printf("read message error: %v", err)
			return
		}
		for _, m := range msg {

			if m.OpCode.IsControl() {

				if m.OpCode == ws.OpPong {
					connection.didReceivePong = true

					if connection.onHb != nil {
						connection.onHb()
					}

					continue
				}

				err := wsutil.HandleClientControlMessage(connection.ws, m)

				if err != nil {
					log.Printf("handle control error: %v", err)
					return
				}

				continue
			}

			if m.OpCode == ws.OpBinary {
				if connection.onMessage != nil {
					connection.onMessage(m.Payload)
				}
			}
		}
	}
}

func (connection *OWebSocketConnection) writeMessages() {
	pingPongTimer := time.NewTimer(pingPeriod)
	waitingForPong := false

	defer func() {
		pingPongTimer.Stop()
		connection.Close()
	}()

	for {

		select {
		case message, ok := <-connection.messageSendChannel:
			{
				if connection.isClosed {
					return
				}

				// If send channel has closed then close websocket connection
				if !ok {

					//_ = connection.write(ws.OpClose, []byte{})
					return
				}

				// Gob WS websockets
				// If an error occurred while sending messages then we have a problem
				if err := wsutil.WriteServerMessage(connection.ws, ws.OpBinary, message); err != nil {
					return
				}
			}

		case <-pingPongTimer.C:
			{
				pingPongTimer.Stop()

				if connection.isClosed {
					return
				}

				if waitingForPong {

					waitingForPong = false

					if connection.didReceivePong {
						connection.didReceivePong = false
						pingPongTimer.Reset(pingPeriod)

					} else {
						connection.Close()
					}

				} else {

					/*
					// If we can't send a ping then it's probably closed
					_, _ = connection.ws.Write(ws.CompiledPing)

					waitingForPong = true
					pingPongTimer.Reset(pongWait)*/
				}
			}
		}
	}
}
