package core

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lisomatrix/channels/channels/auth"

	"github.com/rs/xid"
)

// RemoveIndex - Helper to remove index from slice
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// RemoveChannelIndex - Helper to remove index from slice
func RemoveChannelIndex(s []*HubChannel, index int) []*HubChannel {
	return append(s[:index], s[index+1:]...)
}

// Session - an updated session handling
type Session struct {
	ID                 string
	connection         Connection
	isWaitingForAck    bool
	clientID           string
	identity           *auth.Identity // User AppID and UserID
	deviceID           string         // DeviceID is needed so we the same client can have multiple connections from different devices
	isClosed           bool
	hub                *Hub // Hub for the client AppID
	SubscribedChannels []*HubChannel
	AllowedChannels    []string
	SessionIdentifier  string // We create a string once and store now, instead of creating every time
}

// Init - initialize properties and start sending messages
func (session *Session) Init(connection Connection, deviceID string, identity *auth.Identity, clientID string, hub *Hub) {

	// Fetch channels from cache
	channelIds, found := GetEngine().GetCacheStorage().GetClientChannels(identity.ClientID)

	// If not found fetch from database
	if !found {
		ids, err := GetEngine().GetChannelRepository().GetClientAllowedChannels(identity.ClientID)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "SessionInit: failed to load client allowed channels %v\n", err)
			channelIds = make([]string, 0)
		} else {
			channelIds = ids
		}
	}

	// Initialize session

	if deviceID == "" {
		deviceID = xid.New().String()
	}

	// Device
	session.deviceID = deviceID

	// channels
	session.SubscribedChannels = make([]*HubChannel, 0)
	session.AllowedChannels = channelIds
	// Connection
	session.connection = connection
	// Client info
	session.clientID = clientID
	session.identity = identity
	// App
	session.hub = hub
	// Session ID
	session.ID = xid.New().String()
	session.SessionIdentifier = session.clientID + "-" + session.deviceID

	// Set handlers
	connection.SetOnMessage(session.onNewMessage)
	connection.SetOnClose(session.onClose)
	connection.SetOnHeartBeat(session.onHeartBeat)

	// Update user device online status
	GetEngine().GetPresence().UpdateClientTimestamp(session.clientID)
}

// AddChannel - Add channel while client is connected
func (session *Session) AddChannel(channelID string) {
	session.AllowedChannels = append(session.AllowedChannels, channelID)

	newEvent := NewEvent{
		Type:    NewEvent_NEW_CHANNEL,
		Payload: []byte(channelID),
	}

	data, err := newEvent.Marshal()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Session Add channel: failed to marhal new event: %v\n", err)
	}

	session.connection.Send(data)
}

// RemoveChannel - Remove channel while client is connected
func (session *Session) RemoveChannel(channelID string) {

	newEvent := NewEvent{
		Type:    NewEvent_REMOVE_CHANNEL,
		Payload: []byte(channelID),
	}

	data, err := newEvent.Marshal()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Session Remove channel: failed to marhal new event: %v\n", err)
	}

	session.connection.Send(data)

	var found = false

	for index, channel := range session.SubscribedChannels {
		if channel == channel {

			session.SubscribedChannels = RemoveChannelIndex(session.SubscribedChannels, index)
			found = true
			break
		}
	}

	if found {
		for index, channel := range session.AllowedChannels {
			if channel == channelID {
				session.hub.Unsubscribe(channel, session)

				session.AllowedChannels = RemoveIndex(session.AllowedChannels, index)
				return
			}
		}
	}
}

// Publish - Send data to subscribed client
func (session *Session) Publish(data []byte) {

	if session.isClosed {
		return
	}

	//session.connection.Send(data)
	session.connection.Send(data)
}

func (session *Session) onHeartBeat() {
	// Update timestamps
	if session.isClosed {
		return
	}

	// Update device timestamp
	GetEngine().GetPresence().UpdateClientTimestamp(session.clientID)
}

func (session *Session) onClose() {
	session.Close()
}

func (session *Session) onNewMessage(data []byte) {
	var newEvent NewEvent

	err := newEvent.Unmarshal(data)

	if err != nil {
		log.Println(err)
		return
	}

	if newEvent.Type == NewEvent_SUBSCRIBE {

		var channelSub SubscribeRequest

		err = channelSub.Unmarshal(newEvent.Payload)

		if err != nil {
			log.Println(err)
			return
		}

		didSubscribe := session.CanSubscribe(channelSub.ChannelID)

		session.notifyAck(channelSub.ID, didSubscribe)

	} else if newEvent.Type == NewEvent_PUBLISH {

		var channelPubRequest PublishRequest

		err := channelPubRequest.Unmarshal(newEvent.Payload)

		if err != nil {
			log.Println(err)
			return
		}

		var channelEvent = ChannelEvent{
			SenderID:  session.identity.ClientID,
			EventType: channelPubRequest.EventType,
			Payload:   channelPubRequest.Payload,
			ChannelID: channelPubRequest.ChannelID,
			Timestamp: time.Now().Unix(),
		}

		session.CanPublish(channelPubRequest.ChannelID, &channelEvent, &channelPubRequest)
	}

}

// notifyAck - Notify publish success
func (session *Session) notifyAck(requestID uint32, status bool) {
	ack := PublishAck{
		ReplyTo: requestID,
		Status:  status,
	}

	data, err := ack.Marshal()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Session Notify: failed to marhal ack: %v\n", err)
		return
	}

	newEvent := NewEvent{
		Type:    NewEvent_ACK,
		Payload: data,
	}

	data, err = newEvent.Marshal()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Session Notify: failed to marhal event: %v\n", err)
		return
	}

	session.connection.Send(data)
}

// CanPublish - Check if user is allowed to publish, if so publish
// Also, if a requestID is given we notify the channel (if it is persistent) to store the event
// Otherwise we publish but won't store the event, nor send the notify back
func (session *Session) CanPublish(channelID string, event *ChannelEvent, publishRequest *PublishRequest) {

	for _, c := range session.AllowedChannels {
		if c == channelID {

			didPublish := session.hub.Publish(channelID, event, publishRequest.ID != 0)

			//* INFO: If ID == 0 then we don't need a response back and it won't be stored
			if publishRequest != nil && publishRequest.ID != 0 {
				session.notifyAck(publishRequest.ID, didPublish)
			}

			return
		}
	}
}

// GetIdentifier - Get client and device identifier
func (session *Session) GetIdentifier() string {
	return session.SessionIdentifier
}

// CanSubscribe - Check if user is allowed to subscribe, if so subscribe
func (session *Session) CanSubscribe(channelID string) bool {

	if session.identity.IsAdminKind() {
		channel := session.hub.Subscribe(channelID, session)

		if channel == nil {
			return false
		}

		session.SubscribedChannels = append(session.SubscribedChannels, channel)

		return true
	}

	for _, c := range session.AllowedChannels {
		if c == channelID {
			channel := session.hub.Subscribe(channelID, session)

			if channel == nil {
				return false
			}

			session.SubscribedChannels = append(session.SubscribedChannels, channel)

			return true
		}
	}

	return false
}

// Close - closes session and connection
func (session *Session) Close() {
	session.isClosed = true

	if session.connection.IsConnected() {
		session.connection.Close()
	}

	//GetEngine().GetPresence().SetDeviceOffline(session.clientID, session.deviceID)

	session.hub.RemoveClient(session)
}
