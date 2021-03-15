package connection

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
	//maxMessageSize = 4500
)

type wsMsg struct {
	dataType int
	payload  []byte
}

// WebSocketConnection - Implementation of the Connection interface for WebSocket clients
type WebSocketConnection struct {
	messageSendChannel chan wsMsg
	ws                 *websocket.Conn
	onCloseCB          func()
	onMessage          func([]byte)
	isClosed           bool
}

// IsConnected - If connection is still alive
func (connection *WebSocketConnection) IsConnected() bool {
	return !connection.isClosed
}

// Init - Initialize connection and start reading and writing messages
func (connection *WebSocketConnection) Init(ws *websocket.Conn) {

	connection.ws = ws
	connection.messageSendChannel = make(chan wsMsg, 5)

	go connection.readMessages()
	go connection.writeMessages()
}

// Send - Send the data to the client
func (connection *WebSocketConnection) Send(data []byte) {
	connection.messageSendChannel <- wsMsg{
		dataType: websocket.BinaryMessage,
		payload:  data,
	}
}

// SendText - Send the data to the client
func (connection *WebSocketConnection) SendText(data []byte) {
	connection.messageSendChannel <- wsMsg{
		dataType: websocket.TextMessage,
		payload:  data,
	}
}

// SetOnMessage - Set user sent message callback
func (connection *WebSocketConnection) SetOnMessage(cb func([]byte)) {
	connection.onMessage = cb
}

// SetOnClose - Set on closed connection handler
func (connection *WebSocketConnection) SetOnClose(cb func()) {
	connection.onCloseCB = cb
}

// Close - Close the current connection
func (connection *WebSocketConnection) Close() {
	_ = connection.ws.Close()

	connection.isClosed = true

	if connection.onCloseCB != nil {
		connection.onCloseCB()
	}
}

func (connection *WebSocketConnection) write(mt int, payload []byte) error {
	_ = connection.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return connection.ws.WriteMessage(mt, payload)
}

func (connection *WebSocketConnection) readMessages() {
	defer func() {
		connection.Close()
	}()

	connection.ws.SetReadLimit(maxMessageSize)
	_ = connection.ws.SetReadDeadline(time.Now().Add(pongWait))

	connection.ws.SetPongHandler(func(string) error {
		_ = connection.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := connection.ws.ReadMessage()

		if connection.isClosed {
			return
		}

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			log.Printf("error: %v", err)
			break
		}

		if connection.onMessage != nil {
			go connection.onMessage(msg)
		}
	}
}

var (
	newline = []byte{'\n'}
)

func (connection *WebSocketConnection) writeMessages() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
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
					if err := connection.write(websocket.CloseMessage, []byte{}); err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "Error writing close message: %v\n", err)
					}

					return
				}

				w, err := connection.ws.NextWriter(websocket.TextMessage)

				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Preparing to write message: %v\n", err)
					return
				}

				if _, err = w.Write(message.payload); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Error writing payload: %v\n", err)
				}

				n := len(connection.messageSendChannel)
				for i := 0; i < n; i++ {
					_, _ = w.Write(newline)
					_, err = w.Write((<-connection.messageSendChannel).payload)

					if err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "Error writing payload: %v\n", err)
					}
				}

				if err := w.Close(); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Error closing connection: %v\n", err)
					return
				}

				/*
					// If an error occurred while sending messages then we have a problem
					if err := connection.write(message.dataType, message.payload); err != nil {
						return
					}*/
			}

		case <-ticker.C:
			{
				if connection.isClosed {
					return
				}

				err := connection.ws.SetWriteDeadline(time.Now().Add(writeWait))

				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Error setting ping deadline: %v\n", err)
				}

				// If we can't send a ping then it probably is closed
				if err := connection.write(websocket.PingMessage, []byte{}); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Error writing ping message: %v\n", err)
					return
				}
			}

		}
	}
}
