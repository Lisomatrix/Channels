package core

// PushRequestItem - Push notification request item
type PushRequestItem struct {
	ChannelID string
	EventType string
	Payload   string
	Timestamp int64
	ClientIDs []string
}

type PushNotificationHandler interface {
	EnqueueRequest(request *PushRequestItem)
}