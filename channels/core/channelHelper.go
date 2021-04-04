package core

import (
	"fmt"
	"log"
	"os"
	"time"
)

func SendPushNotification(appID string, channelEvent *ChannelEvent) bool {
	clientIDs, err := GetEngine().GetChannelRepository().GetChannelClients(appID, channelEvent.ChannelID)

	if err != nil {
		log.Println(err)
		return false
	}

	if clientIDs == nil {
		return false
	}

	GetEngine().GetPushHandler().EnqueueRequest(&PushRequestItem{
		ChannelID: channelEvent.ChannelID,
		EventType: channelEvent.EventType,
		Payload:   channelEvent.Payload,
		Timestamp: channelEvent.Timestamp,
		ClientIDs: clientIDs,
	})

	return true
}

// CreateChannel - Validates input an tries to create a channel
func CreateChannel(appID string, channel *Channel) (bool, error) {

	if ChannelExists(appID, channel.ID) {
		return false, nil
	}

	err := GetEngine().GetChannelRepository().CreateChannel(
		channel.ID,
		appID,
		channel.Name,
		channel.CreatedAt,
		false,
		channel.Extra,
		channel.Persistent,
		channel.Private,
		channel.Presence,
		channel.Push,
	)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "HTTP Create Channel failed %v\n", err)
		return false, nil
	}

	// Store new channel in cache
	GetEngine().GetCacheStorage().StoreChannel(appID, channel.ID, channel)

	return true, nil
}

// ChannelExists - Check if channel exists in cache or database
func ChannelExists(appID string, channelID string) bool {

	exists := GetEngine().GetCacheStorage().CheckChannelExistence(appID, channelID)

	if !exists {
		// Check if channel already exists
		if existant, _ := GetEngine().GetChannelRepository().GetAppChannel(appID, channelID); existant != nil {

			// Update cache
			GetEngine().GetCacheStorage().StoreChannel(appID, channelID, existant)

			return true
		}
	}

	return exists
}

// GetChannel - Get channel first from cache, then retry on database and update cache
func GetChannel(appID string, channelID string) (*Channel, error) {

	channel := GetEngine().GetCacheStorage().GetChannel(appID, channelID)

	if channel != nil {
		return channel, nil
	}

	channel, err := GetEngine().GetChannelRepository().GetAppChannel(appID, channelID)

	if err != nil {

		fmt.Fprintf(os.Stderr, "Get channel: failed to get app channel %v\n", err)
		return nil, err
	}

	if channel == nil {
		return nil, nil
	}

	// Update cache
	GetEngine().GetCacheStorage().StoreChannel(appID, channelID, channel)

	return channel, nil
}

// JoinChannel - Join client to a given channel, and update cache and current connected and affected clients
func JoinChannel(appID string, channelID string, clientID string) (bool, error) {

	client, err := GetClient(appID, clientID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Join channel: failed to get app client %v\n", err)
		return false, nil
	}

	if client == nil {
		return false, nil
	}

	channel, err := GetChannel(appID, channelID)

	if err != nil {
		return false, err
	}

	if channel == nil {
		return false, nil
	}

	err = GetEngine().GetChannelRepository().JoinClient(appID, channelID, clientID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Join channel: failed to join client to channel %v\n", err)
		return false, err
	}

	// Update cache
	GetEngine().GetCacheStorage().AddClientChannel(clientID, channelID)

	// Notify current connected clients
	hub := GetEngine().HubsHandler.ContainsHub(appID)

	if hub != nil {
		hub.AddChannelToClient(clientID, channelID)
	}

	// If the channel cares about persistence and presence
	// Then store this new event
	if channel.Persistent && channel.Presence {
		newChannelEvent := &ChannelEvent{
			SenderID:  clientID,
			Timestamp: time.Now().Unix(),
			EventType: "Join",
			ChannelID: channelID,
		}

		clientJoined := ClientJoin{
			ChannelID:            channelID,
			ClientID:             clientID,
		}

		// Store and cache new event
		GetEngine().StoreEvent(channel.AppID, newChannelEvent)
		GetEngine().GetCacheStorage().StoreChannelEvent(channelID, appID, newChannelEvent)

		// If there are clients connected to hub and channel
		// Then publish to them, otherwise there is no point
		// But we still need to publish to other servers
		if hub != nil {
			channel := hub.ContainsChannel(channelID)

			if channel != nil {

				if data, err := clientJoined.Marshal(); err == nil {
					// Publish to local clients only, we send to other servers after
					channel.PublishJoinLeave(NewEvent_JOIN_CHANNEL, data)
				} else {
					_, _ = fmt.Fprintf(os.Stderr, "Join channel: failed to marshal join client event %v\n", err)
				}
			}

		}

		// Then we publish to other servers
		//GetEngine().GetPublisher().PublishChannelEvent(appID, channelID, newChannelEvent)
	}

	if channel.Presence {
		// Publish presence event to other servers
		GetEngine().GetPublisher().PublishChannelPresenceChange(appID, channelID, clientID, true)
	}

	// Notify clientID in other servers that he received access to channel
	GetEngine().GetPublisher().PublishChannelAccessChange(appID, channelID, clientID, true)

	return true, err
}

// LeaveChannel - Remove client from a given channel, and update cache
func LeaveChannel(appID string, channelID string, clientID string) (bool, error) {

	client, err := GetClient(appID, clientID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Leave channel: failed to get app client %v\n", err)
		return false, nil
	}

	if client == nil {
		return false, nil
	}

	channel, err := GetChannel(appID, channelID)

	if err != nil {
		return false, err
	}

	if channel == nil {
		return false, nil
	}

	// Remove client from channel
	err = GetEngine().GetChannelRepository().LeaveClient(appID, channelID, clientID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Leave channel: failed to remove client from channel %v\n", err)
		return false, err
	}

	// Update cache
	GetEngine().GetCacheStorage().RemoveClientChannel(clientID, channelID)

	// Notify current connected clients
	hub := GetEngine().HubsHandler.ContainsHub(appID)

	if hub != nil {
		hub.RemoveChannelFromClient(clientID, channelID)
	}

	// If the channel cares about persistence and presence
	// Then store this new event
	if channel.Persistent && channel.Presence {
		newChannelEvent := &ChannelEvent{
			SenderID:  clientID,
			Timestamp: time.Now().Unix(),
			EventType: "Leave",
			ChannelID: channelID,
			Payload:   clientID,
		}

		// Store and cache new event
		GetEngine().StoreEvent(channel.AppID, newChannelEvent)
		GetEngine().GetCacheStorage().StoreChannelEvent(channelID, appID, newChannelEvent)

		clientLeave := &ClientLeave{
			ChannelID:            channelID,
			ClientID:             clientID,
		}

		// If there are clients connected to hub and channel
		// Then publish to them, otherwise there is no point
		// But we still need to publish to other servers
		if hub != nil {
			channel := hub.ContainsChannel(channelID)

			if data, err := clientLeave.Marshal(); err == nil {
				// Publish to local clients only, we send to other servers after
				channel.PublishJoinLeave(NewEvent_LEAVE_CHANNEL, data)
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "Leave channel: failed to marshal leave client event %v\n", err)
			}

		}

		// Then we publish to other servers
		//GetEngine().GetPublisher().PublishChannelEvent(appID, channelID, newChannelEvent)
	}

	if channel.Presence {
		// Publish presence event to other servers
		GetEngine().GetPublisher().PublishChannelPresenceChange(appID, channelID, clientID, false)
	}

	// Notify clientID in other servers that he lost access to channel
	GetEngine().GetPublisher().PublishChannelAccessChange(appID, channelID, clientID, false)

	return true, err
}

// DeleteChannel - Delete channel from database and cache, and notify all connected clients
func DeleteChannel(appID string, channelID string) (bool, error) {

	channel, err := GetChannel(appID, channelID)

	if err != nil {
		return false, err
	}

	if channel == nil {
		return false, nil
	}

	if err := GetEngine().GetChannelRepository().DeleteChannel(appID, channelID); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "HTTP Delete channel: failed to delete channel %v\n", err)
		return false, err
	}

	// Update cache
	GetEngine().GetCacheStorage().RemoveChannel(appID, channelID)

	// Notify current connected clients
	GetEngine().HubsHandler.GetHub(appID).DeleteChannel(channelID)

	return true, nil
}

// SetChannelCloseStatus - Set channel close status
func SetChannelCloseStatus(appID string, channelID string, closed bool) (bool, error) {

	channel, err := GetChannel(appID, channelID)

	if err != nil {
		return false, err
	}

	if channel == nil {
		return false, nil
	}

	channel.IsClosed = closed

	if err := GetEngine().GetChannelRepository().SetChannelCloseStatus(appID, channelID, closed); err != nil {
		fmt.Fprintf(os.Stderr, "Set channel close status: failed to save %v\n", err)
		return false, err
	}

	if closed {
		// Notify all connected client in case it is closing channel
		GetEngine().HubsHandler.GetHub(appID).DeleteChannel(channelID)
		// Remove from cache
		GetEngine().GetCacheStorage().RemoveChannel(appID, channelID)
	} else {
		// Update cache
		GetEngine().GetCacheStorage().StoreChannel(appID, channelID, channel)
	}

	return true, nil
}
