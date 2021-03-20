package core

import (
	"strings"
	"sync"
)

// NewHub - Create a new Hub
func NewHub(AppID string) *Hub {
	return &Hub{
		AppID: AppID,
	}
}

// Hub - Handles channels and publishing
type Hub struct {
	AppID            string
	channels         sync.Map //[string]*Channel
	connectedClients sync.Map //[string]*Session
}

// DeleteChannel - Remove channel including subscriptions
func (hub *Hub) DeleteChannel(channelID string) {
	value, loaded := hub.channels.LoadAndDelete(channelID)

	if loaded {
		channel := value.(*HubChannel)
		channel.DeleteChannel()
	}
}

// AddChannelToClient - Add channel to current connected client
func (hub *Hub) AddChannelToClient(clientID string, channelID string) {
	hub.connectedClients.Range(func(key interface{}, value interface{}) bool {

		id := key.(string)

		if strings.Contains(id, clientID) {
			session := value.(*Session)

			session.AddChannel(channelID)
		}

		return true
	})
}

// RemoveChannelFromClient - Remove channel to current connected client
func (hub *Hub) RemoveChannelFromClient(clientID string, channelID string) {
	hub.connectedClients.Range(func(key interface{}, value interface{}) bool {

		id := key.(string)

		if strings.Contains(id, clientID) {
			session := value.(*Session)

			session.RemoveChannel(channelID)
		}

		return true
	})
}

// AddClient - Add client to connected map
func (hub *Hub) AddClient(session *Session) {
	hub.connectedClients.Store(session.GetIdentifier(), session)
}

// RemoveClient - Remove client from connected clients and channels
func (hub *Hub) RemoveClient(session *Session) {
	_, isOK := hub.connectedClients.LoadAndDelete(session.GetIdentifier())

	if !isOK {
		return
	}

	for _, channel := range session.SubscribedChannels {
		hub.removeSessionFromChannel(channel.Data.ID, session)
	}
}

func (hub *Hub) removeSessionFromChannel(channelID string, session *Session) {
	data, isOK := hub.channels.Load(channelID)

	if !isOK {
		return
	}

	chann := data.(*HubChannel)
	chann.RemoveClient(session)
}

// Close - Remove all channels and connections
func (hub *Hub) Close() {

}

// AddChannel - Add channel to hub
func (hub *Hub) AddChannel(id string) {

	chann := NewChannel(id, hub.AppID, hub)

	if chann == nil {
		return
	}

	hub.channels.Store(id, chann)
}

// ContainsChannel - Get HubChannel if exists in memory
func (hub *Hub) ContainsChannel(channelID string) *HubChannel {
	// Get channel from local cache
	data, isOk := hub.channels.Load(channelID)

	if isOk {
		return data.(*HubChannel)
	}

	return nil
}

// Publish - Send the given payload to subscribed session
func (hub *Hub) Publish(channelID string, channelEvent *ChannelEvent, shouldStore bool) bool {

	// Get channel from local cache
	data, isOk := hub.channels.Load(channelID)

	var chann *HubChannel

	if !isOk {

		// Load channel
		chann = NewChannel(channelID, hub.AppID, hub)

		// If not found cancel publish
		if chann == nil {
			return false
		}

		// If found cache it
		hub.channels.Store(channelID, chann)

	} else {
		// Found in cache
		chann = data.(*HubChannel)
	}

	// Publish event
	chann.Publish(channelEvent, shouldStore)

	return true
}

// Subscribe - Add subscriber to given channel
func (hub *Hub) Subscribe(channelID string, session *Session) *HubChannel {
	data, isOK := hub.channels.Load(channelID)

	var chann *HubChannel

	// If there are other subscribers already then add this one
	if isOK {
		chann = data.(*HubChannel)
	} else {
		chann = NewChannel(channelID, hub.AppID, hub)

		// If is nil then there are no channels created
		if chann == nil {
			return nil
		}
	}

	chann.NewClient(session)
	hub.channels.Store(channelID, chann)
	return chann
}

// Unsubscribe - Remove subscriber from given channel
func (hub *Hub) Unsubscribe(channelID string, session *Session) {
	data, isOK := hub.channels.Load(channelID)

	if !isOK {
		return
	}

	var chann = data.(*HubChannel)

	chann.RemoveClient(session)
}
