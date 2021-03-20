package core

// PublishHandler - Interface for publishing events between servers
type PublishHandler interface {
	PublishChannelPresenceChange(appID string, channelID string, clientID string, isJoin bool)
	PublishChannelAccessChange(appID string, channelID string, clientID string, isAdd bool)
	PublishChannelEvent(appID string, channelID string, channelEvent *ChannelEvent)
	PublishChannelOnlineChange(appID string, channelID string, statusUpdate *OnlineStatusUpdate)
	Subscribe(appID string, channelID string)
	Unsubscribe(appID string, channelID string)
}
