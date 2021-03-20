package core

import (
	"fmt"
	"go.uber.org/atomic"
	"os"
	"sync"
	"time"
)

// ChannelPublish - Client publish payload
type ChannelPublish struct {
	ChannelID string
	Payload   string
	EventType string
}

// ChannelsOptions - Channel options
type ChannelsOptions struct {
	Persistent bool
}

// NewChannel - Create and initialize channel
func NewChannel(ID string, AppID string, hub *Hub) *HubChannel {

	chann := GetEngine().GetCacheStorage().GetChannel(AppID, ID)

	if chann == nil {
		c, err := GetEngine().GetChannelRepository().GetAppChannel(AppID, ID)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "NewChannel: failed to fetch channel: %v\n", err)
			return nil
		}

		chann = c

		if chann == nil {
			_, _ = fmt.Fprintf(os.Stderr, "NewChannel: attempting to create unexistent channel")
			return nil
		}

		// Update cache
		GetEngine().GetCacheStorage().StoreChannel(AppID, ID, c)
	}

	hubChannel := &HubChannel{
		Data: chann,
		hub: hub,
	}

	if chann.Presence {

		clients, err := GetEngine().GetChannelRepository().GetChannelClients(chann.AppID, chann.ID)

		var zero int64 = 0

		if err == nil {
			for _, c := range clients {
				clientStatus := ClientStatus{
					Status:    false,
					Timestamp: zero,
				}
				hubChannel.connectedClientsStatus.Store(c, clientStatus)
			}
		}

		presences := GetEngine().GetPresence().GetChannelClientsPresence(chann.AppID, chann.ID)

		if presences != nil {

			for key, value := range presences {
				clientStatus := ClientStatus{
					Status:    false,
					Timestamp: value,
				}
				hubChannel.connectedClientsStatus.Store(key, clientStatus)
			}

		}

	}

	GetEngine().GetPublisher().Subscribe(chann.AppID, chann.ID)

	return hubChannel
}

// HubChannel - Handler for topic
type HubChannel struct {
	Data                   *Channel
	connectedUsers         sync.Map //[string(session_identifier)]*Session
	connectedClientsStatus sync.Map //[string(clientID)]TimeStamp
	isClosing              bool
	connectedCounter	   atomic.Int32
	hub 				   *Hub
}

// DeleteChannel - Unsubscribe all clients and stop accepting subscriptions
func (channel *HubChannel) DeleteChannel() {

	channel.isClosing = true

	channel.connectedUsers.Range(func(key interface{}, value interface{}) bool {

		session := value.(*Session)

		session.RemoveChannel(channel.Data.ID)

		GetEngine().GetPublisher().Unsubscribe(channel.Data.ID, channel.Data.AppID)

		return true
	})
}

// ExternalPublish - Publish to be used by HTTP and Publisher so we don't republish nor store in db/cache
func (channel *HubChannel) ExternalPublish(channelEvent *ChannelEvent) bool {
	if channel.isClosing {
		return false
	}

	// * We parse the message here, so
	// * we avoid parsing for each connection
	//data, err := json.Marshal(channelEvent)
	data, err := channelEvent.Marshal()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Session Publish: failed to marhal channel event: %v\n", err)
		return false
	}

	channel.connectedUsers.Range(func(key interface{}, value interface{}) bool {

		session := value.(*Session)

		session.Publish(data)

		return true
	})

	return true
}

// Publish - Send message to all connected clients
func (channel *HubChannel) Publish(channelEvent *ChannelEvent, shouldStore bool) bool {

	if channel.isClosing {
		return false
	}

	// * We parse the message here, so
	// * we avoid parsing for each connection
	//data, err := json.Marshal(channelEvent)

	data, err := channelEvent.Marshal()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Session Publish: failed to marhal channel event: %v\n", err)
		return false
	}

	newEvent := NewEvent{
		Type:    NewEvent_PUBLISH,
		Payload: data,
	}

	newEventData, err := newEvent.Marshal()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Session Publish: failed to marhal new channel event: %v\n", err)
		return false
	}

	// If it is a persistent channel store message in DB and cache
	if channel.Data.Persistent && shouldStore {
		GetEngine().StoreEvent(channel.Data.AppID, channelEvent)
		GetEngine().GetCacheStorage().StoreChannelEvent(channel.Data.ID, channelEvent)
	}

	GetEngine().GetPublisher().PublishChannelEvent(channel.Data.AppID, channel.Data.ID, channelEvent)

	channel.connectedUsers.Range(func(key interface{}, value interface{}) bool {

		session := value.(*Session)

		session.Publish(newEventData)

		return true
	})

	return true
}

// ExternalPublishStatusChange - Publish new event about user status update, it doesn't resend data back to publisher
func (channel *HubChannel) ExternalPublishStatusChange(statusUpdate *OnlineStatusUpdate) bool {
	if channel.isClosing {
		return false
	}

	// Update channel status
	clientStatus := ClientStatus{
		Status:    statusUpdate.Status,
		Timestamp: statusUpdate.Timestamp,
	}
	channel.connectedClientsStatus.Store(statusUpdate.ClientID, clientStatus)

	// * We parse the message here, so
	// * we avoid parsing for each connection
	statusUpdateData, err := statusUpdate.Marshal()
	//statusUpdateData, err := json.Marshal(statusUpdate)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Session Publish: failed to marhal status update: %v\n", err)
	}

	newEvent := NewEvent{
		Type:    NewEvent_ONLINE_STATUS,
		Payload: statusUpdateData,
	}

	//newEventData, err := json.Marshal(&newEvent)
	newEventData, err := newEvent.Marshal()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Session Publish: failed to marhal new event on status update: %v\n", err)
	}

	channel.connectedUsers.Range(func(key interface{}, value interface{}) bool {

		session := value.(*Session)

		// Update other users only
		if session.clientID != statusUpdate.ClientID {
			session.Publish(newEventData)
		}

		return true
	})

	return true
}

// PublishStatusChange - Publish new event about user status update
func (channel *HubChannel) PublishStatusChange(statusUpdate *OnlineStatusUpdate) bool {
	if channel.isClosing {
		return false
	}

	clientStatus := ClientStatus{
		Status:    statusUpdate.Status,
		Timestamp: statusUpdate.Timestamp,
	}
	// Update channel status
	channel.connectedClientsStatus.Store(statusUpdate.ClientID, clientStatus)

	// * We parse the message here, so
	// * we avoid parsing for each connection
	//statusUpdateData, err := json.Marshal(statusUpdate)
	statusUpdateData, err := statusUpdate.Marshal()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Session Publish: failed to marhal status update: %v\n", err)
	}

	newEvent := NewEvent{
		Type:    NewEvent_ONLINE_STATUS,
		Payload: statusUpdateData,
	}

	//newEventData, err := json.Marshal(&newEvent)
	newEventData, err := newEvent.Marshal()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Session Publish: failed to marhal new event on status update: %v\n", err)
	}

	// Update other servers about this change
	GetEngine().GetPublisher().PublishChannelOnlineChange(channel.Data.AppID, channel.Data.ID, statusUpdate)

	channel.connectedUsers.Range(func(key interface{}, value interface{}) bool {

		session := value.(*Session)

		// Update other users only
		if session.clientID != statusUpdate.ClientID {
			session.Publish(newEventData)
		}

		return true
	})

	return true
}

// NewClient - Add client to channel
func (channel *HubChannel) NewClient(session *Session) {

	if channel.isClosing {
		return
	}

	// Add connected counter
	channel.connectedCounter.Inc()

	channel.connectedUsers.Store(session.GetIdentifier(), session)

	if channel.Data.Presence {
		channel.shouldNotifyOnlinePresenceChange(session)
		// Prepare initial state
		initialPresenceState := InitialPresenceStatus{
			ChannelID:    channel.Data.ID,
			ClientStatus: make(map[string]*ClientStatus),
		}

		// Use channel local state
		channel.connectedClientsStatus.Range(func(key interface{}, value interface{}) bool {

			clientID := key.(string)

			if clientID == session.clientID {
				return true
			}

			status := value.(ClientStatus)

			initialPresenceState.ClientStatus[clientID] = &status

			return true
		})

		// Prepare and marshal response
		//data, err := json.Marshal(&initialPresenceState)
		data, err := initialPresenceState.Marshal()

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Channel Initial presence: failed to marshal initial presence status %v\n", err)
			return
		}

		newEvent := NewEvent{
			Type:    NewEvent_INITIAL_ONLINE_STATUS,
			Payload: data,
		}

		//eventData, err := json.Marshal(&newEvent)
		eventData, err := newEvent.Marshal()

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Channel Initial presence: failed to marshal initial presence status of new event %v\n", err)
			return
		}

		// Send state
		session.Publish(eventData)
	}
}

func (channel *HubChannel) shouldNotifyOnlinePresenceChange(session *Session) {

	timeStamp := time.Now().Unix()

	status := ClientStatus{
		Status:    true,
		Timestamp: timeStamp,
	}

	// Update channel status
	channel.connectedClientsStatus.Store(session.clientID, status)

	GetEngine().GetPresence().AddOnlineChannelDevice(channel.Data.AppID, channel.Data.ID, session.clientID, session.deviceID)

	// Get how many are left
	amount := GetEngine().GetPresence().GetChannelAmountOfClientDevices(channel.Data.AppID, channel.Data.ID, session.clientID)

	if amount == 1 {
		statusUpdate := OnlineStatusUpdate{
			ChannelID: channel.Data.ID,
			ClientID:  session.clientID,
			Status:    true, // If not remove is online
			Timestamp: timeStamp,
		}

		channel.PublishStatusChange(&statusUpdate)
	}
}

// shouldNotifyOfflinePresenceChange - Check if an offline status update should be done
func (channel *HubChannel) shouldNotifyOfflinePresenceChange(session *Session) {

	// Remove this devices from channel online devices
	GetEngine().GetPresence().RemoveOnlineChannelDevice(channel.Data.AppID, channel.Data.ID, session.clientID, session.deviceID)

	// Set a timer of X seconds
	// To prevent device on reconnecting to constantly change online status
	go func() {

		timer := time.NewTimer(time.Second * 15)

		// wait for timer
		<-timer.C

		// Check if the remove device is connected
		// If so then it reconnected and there is no need to publish the status update
		if GetEngine().GetPresence().IsClientDeviceConnectToChannel(channel.Data.AppID, channel.Data.ID, session.clientID, session.deviceID) {
			timer.Stop()
			return
		}

		// If the device is not back online
		// We must check if client is not connected with another device
		// If the he is, then ignore the status update
		amount := GetEngine().GetPresence().GetChannelAmountOfClientDevices(channel.Data.AppID, channel.Data.ID, session.clientID)

		if amount < 0 {
			amount = 0
		}

		if amount > 0 {
			timer.Stop()
			return
		}

		statusUpdate := OnlineStatusUpdate{
			ChannelID: channel.Data.ID,
			ClientID:  session.clientID,
			Status:    false, // If not remove is online
			Timestamp: time.Now().Unix(),
		}

		status := ClientStatus{
			Status:    statusUpdate.Status,
			Timestamp: statusUpdate.Timestamp,
		}

		// Update channel status
		channel.connectedClientsStatus.Store(session.clientID, status)

		channel.PublishStatusChange(&statusUpdate)
		timer.Stop()
	}()
}

// RemoveClient - Remove client from channel
func (channel *HubChannel) RemoveClient(session *Session) {

	if channel.isClosing {
		return
	}

	channel.connectedCounter.Dec()

	channel.connectedUsers.LoadAndDelete(session.GetIdentifier())

	if channel.Data.Presence {
		channel.shouldNotifyOfflinePresenceChange(session)
	}

	// If we have 0 users, then wait some minutes before removing
	// The channel from the Hub, so we save resources
	if channel.connectedCounter.Load() == 0 {
		channel.shouldCloseChannel()
	}
}

// shouldCloseChannel - Run a timer for 15 minutes, if no new connection shows up, the channel is closed
func (channel *HubChannel) shouldCloseChannel() {
	// Set a timer of X seconds
	// To prevent removing channel if a user connects in the next 15 mins
	go func() {
		timer := time.NewTimer(time.Minute * 15)

		// wait for timer
		<-timer.C
		timer.Stop()

		if channel.connectedCounter.Load() == 0 {
			fmt.Printf("No subscribers on channel %s for the last 15 mins, closing channel...", channel.Data.ID)
			channel.hub.DeleteChannel(channel.Data.ID)
		}

	}()
}