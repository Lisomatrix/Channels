package connection

import (
	"log"
	"net/http"

	"github.com/channelserver/channelserver/auth"
	"github.com/channelserver/channelserver/core"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	ws "github.com/gobwas/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// RequestHandler - Default WebSocket handler
func RequestHandler(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Get token and AppID headers
	token := request.Header.Get("Authorization")
	appID := request.Header.Get("AppID")
	deviceID := request.Header.Get("DeviceID")

	queryValues := request.URL.Query()

	// If the websocket is from the browser then they come as query params
	if token == "" {
		token = queryValues.Get("Authorization")
	}

	if appID == "" {
		appID = queryValues.Get("AppID")
	}

	if deviceID == "" {
		deviceID = queryValues.Get("DeviceID")
	}

	// If in neither the headers or query params it was found
	// Then the request is malformed
	if appID == "" || token == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	identity, isOK := auth.VerifyToken(token)

	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(writer, request, nil)
	//conn, _, _, err := ws.UpgradeHTTP(request, writer)

	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	// Start session handler

	var connection *WebSocketConnection = new(WebSocketConnection)
	var session *core.Session = new(core.Session)

	/*
		client := core.GetEngine().GetCacheStorage().GetClient(appID, identity.ClientID)

		if client == nil {
			c, err := core.GetEngine().GetClientRepository().GetAppClient(appID, identity.ClientID)

			if err != nil {
				fmt.Fprintf(os.Stderr, "WS Handler: Failed to get client data: %v\n", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}

			client = c
		}*/

	hub := core.GetEngine().HubsHandler.GetHub(identity.AppID)

	connection.Init(conn)
	session.Init(connection, deviceID, &identity, identity.ClientID, hub)

	hub.AddClient(session)
}

// OptimizedRequestHandler - Optimized version of WebSocket handshake
func OptimizedRequestHandler(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Get token and AppID headers
	token := request.Header.Get("Authorization")
	appID := request.Header.Get("AppID")
	deviceID := request.Header.Get("DeviceID")

	queryValues := request.URL.Query()

	// If the websocket is from the browser then they come as query params
	if token == "" {
		token = queryValues.Get("Authorization")
	}

	if appID == "" {
		appID = queryValues.Get("AppID")
	}

	if deviceID == "" {
		deviceID = queryValues.Get("DeviceID")
	}

	// If in neither the headers or query params it was found
	// Then the request is malformed
	if appID == "" || token == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	identity, isOK := auth.VerifyToken(token)

	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Upgrade to WebSocket
	conn, _, _, err := ws.UpgradeHTTP(request, writer)

	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	// Start session handler
	var connection *OWebSocketConnection = new(OWebSocketConnection)
	var session *core.Session = new(core.Session)

	hub := core.GetEngine().HubsHandler.GetHub(identity.AppID)

	connection.Init(conn)
	session.Init(connection, deviceID, &identity, identity.ClientID, hub)

	hub.AddClient(session)
}
