package core

// PublishHandler - Interface for publishing events between servers
type PublishHandler interface {
	PublishChannelEvent(appID string, channelID string, channelEvent *ChannelEvent)
	PublishChannelOnlineChange(appID string, channelID string, statusUpdate *OnlineStatusUpdate)
	Subscribe(appID string, channelID string)
	Unsubscribe(appID string, channelID string)
}
