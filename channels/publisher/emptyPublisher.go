package publisher

import "github.com/lisomatrix/channels/channels/core"

type EmptyPublisher struct {}

func (publisher *EmptyPublisher) PublishChannelEvent(appID string, channelID string, channelEvent *core.ChannelEvent) {

}

func (publisher *EmptyPublisher) PublishChannelOnlineChange(appID string, channelID string, statusUpdate *core.OnlineStatusUpdate) {

}

func (publisher *EmptyPublisher) Subscribe(appID string, channelID string) {

}

func (publisher *EmptyPublisher) Unsubscribe(appID string, channelID string) {

}

func NewEmptyPublisher() *EmptyPublisher {
	return &EmptyPublisher{}
}