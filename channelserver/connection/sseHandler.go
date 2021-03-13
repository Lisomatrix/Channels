package connection

/*
// SseHandler - Handle for SSE subscribe
func SseHandler(context *gin.Context) {
	request := context.Request
	rw := context.Writer

	channel := context.Params.ByName("channel")
	appID := request.Header.Get("AppID")
	userID := request.Header.Get("UserID")

	if channel == "" || appID == "" || userID == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if client supports streaming data
	_, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	// Initialize connection and session handler
	var connection *SSEConnection = new(SSEConnection)
	var session *core.Session = new(core.Session)

	client := core.GetAppStorage().LoadClient(userID)

	identity := &auth.Identity{
		AppID:    appID,
		ClientID: userID,
	}

	//client.accountID = tokenData.accountID
	//client.deviceID = tokenData.deviceID

	hub := core.GetHubsHandler().GetHub(identity.AppID)

	connection.Init()
	session.Init(connection, identity, client, hub)

	hub.AddClient(session)

	// Subscribe to channel
	session.CanSubscribe(channel)

	// Listen for when clients closes
	notify := request.Context().Done()

	go func() {
		<-notify
		session.Close()
	}()

	// Stream data into client
	context.Stream(func(w io.Writer) bool {

		if msg, ok := <-connection.messageSendChannel; ok {
			context.SSEvent("message", msg)

			if !connection.isClosed {
				return true
			}
		}

		return false
	})
}
*/
