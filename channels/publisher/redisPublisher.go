// This package holds the publisher interface implementations
package publisher

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lisomatrix/channels/channels/core"

	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
)

// RedisPublisher - Implementation of PublishHandler interface
type RedisPublisher struct {
	client *redis.Client
	pubsub *redis.PubSub
	ctx    context.Context
}

func (publisher *RedisPublisher) PublishChannelPresenceChange(appID string, channelID string, clientID string, isJoin bool) {

	presenceType := ExternalChannelPresenceType_Join

	if !isJoin {
		presenceType = ExternalChannelPresenceType_Leave
	}

	channelPresenceChange := ExternalJoinLeaveClientEvent{
		ClientID:     clientID,
		ChannelID:    channelID,
		PresenceType: presenceType,
	}

	newEvent := ExternalNewEvent{
		Type:              ExternalNewEventType_ChannelPresence,
		ServerID:          core.GetEngine().GetServerID(),
		ExternalJoinLeave: &channelPresenceChange,
	}

	data, err := newEvent.Marshal()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: failed to marshal event %v\n", err)
		return
	}

	cmd := publisher.client.Publish(publisher.ctx, appID+":"+channelID, data)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: failed to publish event %v\n", err)
		return
	}
}

func (publisher *RedisPublisher) PublishChannelAccessChange(appID string, channelID string, clientID string, isAdd bool) {

	accessType := ExternalChannelAccessType_Add

	if !isAdd {
		accessType = ExternalChannelAccessType_Remove
	}

	accessEvent := ExternalChannelAccessEvent{
		ExternalAccessType: accessType,
		ClientID:           clientID,
		ChannelID:          channelID,
	}

	newEvent := ExternalNewEvent{
		Type:                ExternalNewEventType_ChannelAccess,
		ServerID:            core.GetEngine().GetServerID(),
		ExternalAccessEvent: &accessEvent,
	}

	data, err := newEvent.Marshal()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: failed to marshal event %v\n", err)
		return
	}

	cmd := publisher.client.Publish(publisher.ctx, appID+":"+channelID, data)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: failed to publish event %v\n", err)
		return
	}
}

// PublishChannelOnlineChange - Publish Online status change to other servers
func (publisher *RedisPublisher) PublishChannelOnlineChange(appID string, channelID string, statusUpdate *core.OnlineStatusUpdate) {

	onlineStatusEvent := ExternalOnlineStatusEvent{
		ClientID:  statusUpdate.ClientID,
		Timestamp: statusUpdate.Timestamp,
		Status:    statusUpdate.Status,
	}

	newEvent := ExternalNewEvent{
		Type:                 ExternalNewEventType_OnlineStatus,
		ServerID:             core.GetEngine().GetServerID(),
		ExternalOnlineStatus: &onlineStatusEvent,
	}

	data, err := newEvent.Marshal()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: failed to marshal event %v\n", err)
		return
	}

	cmd := publisher.client.Publish(publisher.ctx, appID+":"+channelID, data)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: failed to publish event %v\n", err)
		return
	}
}

// PublishChannelEvent - Send event for other servers listening for this event
func (publisher *RedisPublisher) PublishChannelEvent(appID string, channelID string, channelEvent *core.ChannelEvent) {

	publishEvent := ExternalPublishEvent{
		SenderID:  channelEvent.SenderID,
		Payload:   channelEvent.Payload,
		Timestamp: channelEvent.Timestamp,
		EventType: channelEvent.EventType,
	}

	newEvent := ExternalNewEvent{
		Type:                 ExternalNewEventType_ChannelEvent,
		ServerID:             core.GetEngine().GetServerID(),
		ExternalPublishEvent: &publishEvent,
	}

	data, err := newEvent.Marshal()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: failed to marshal event %v\n", err)
		return
	}

	cmd := publisher.client.Publish(publisher.ctx, appID+":"+channelID, data)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: failed to publish event %v\n", err)
		return
	}
}

// Unsubscribe - Unsubscribe from a channel in redis
func (publisher *RedisPublisher) Unsubscribe(appID string, channelID string) {
	err := publisher.pubsub.Unsubscribe(publisher.ctx, appID+":"+channelID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: failed to unsubscribe from channel %v\n", err)
	}
}

// Subscribe - Subscribe to a channel in redis
func (publisher *RedisPublisher) Subscribe(appID string, channelID string) {
	err := publisher.pubsub.Subscribe(publisher.ctx, appID+":"+channelID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: failed get subscribe to channel %v\n", err)
	}
}

func (publisher *RedisPublisher) handleSubscribeMessages() {

	defer log.Println("Stopped handling redis messages")

	ch := publisher.pubsub.Channel()

	for {
		data, isActive := <-ch

		if !isActive {
			return
		}

		var newEvent ExternalNewEvent

		err := newEvent.Unmarshal([]byte(data.Payload))

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: failed umarshal external event %v\n", err)
			continue
		}

		// We don't want to listen for our own events
		if newEvent.ServerID == core.GetEngine().GetServerID() {
			continue
		}

		_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: received event with size: %d bytes\n", len([]byte(data.Payload)))

		parts := strings.Split(data.Channel, ":")

		appID := parts[0]
		channelID := parts[1]

		hub := core.GetEngine().GetHubsHandler().ContainsHub(appID)

		// If there is no hub then we don't have clients from the hub
		if hub == nil {
			continue
		}

		// If there is no channel then we don't have clients listening to this channel
		channel := hub.ContainsChannel(channelID)

		if channel == nil {
			continue
		}

		if newEvent.Type == ExternalNewEventType_ChannelEvent {

			event := newEvent.GetExternalPublishEvent()

			channel.ExternalPublish(&core.ChannelEvent{
				SenderID:  event.SenderID,
				Payload:   event.Payload,
				EventType: event.EventType,
				Timestamp: event.Timestamp,
				ChannelID: channelID,
			})

		} else if newEvent.Type == ExternalNewEventType_OnlineStatus {
			event := newEvent.GetExternalOnlineStatus()

			channel.ExternalPublishStatusChange(&core.OnlineStatusUpdate{
				ChannelID: channelID,
				Timestamp: event.Timestamp,
				Status:    event.Status,
				ClientID:  event.ClientID,
			})

		} else if newEvent.Type == ExternalNewEventType_ChannelAccess {
			event := newEvent.GetExternalAccessEvent()

			if event.ExternalAccessType == ExternalChannelAccessType_Add {
				hub.RemoveChannelFromClient(event.ClientID, event.ChannelID)
			} else if event.ExternalAccessType == ExternalChannelAccessType_Remove {
				hub.AddChannelToClient(event.ClientID, event.ChannelID)
			}

		} else if newEvent.Type == ExternalNewEventType_ChannelPresence {

			event := newEvent.GetExternalJoinLeave()

			if event.PresenceType == ExternalChannelPresenceType_Join {

				clientJoined := core.ClientJoin{
					ChannelID: event.ChannelID,
					ClientID:  event.ClientID,
				}

				if data, err := clientJoined.Marshal(); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: failed to marshal external channel presence event (JOIN TYPE) \n")
				} else {
					channel.PublishJoinLeave(core.NewEvent_JOIN_CHANNEL, data)
				}
			} else if event.PresenceType == ExternalChannelPresenceType_Leave {
				clientLeave := core.ClientLeave{
					ChannelID: event.ChannelID,
					ClientID:  event.ClientID,
				}

				if data, err := clientLeave.Marshal(); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: failed to marshal external channel presence event (LEAVE TYPE) \n")
				} else {
					channel.PublishJoinLeave(core.NewEvent_LEAVE_CHANNEL, data)
				}
			}

		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Redis Publisher: received Unknown event type \n")
		}

	}
}

// NewRedisPublisher - Create a new instance of redis publisher
func NewRedisPublisher() *RedisPublisher {
	redisPublisher := new(RedisPublisher)
	redisPublisher.ctx = context.Background()

	client := redis.NewClient(
		&redis.Options{
			Addr: "127.0.0.1:6379",
			// Password: "CulP3gnpgSAxFlbjO/JrNCR/uTKFKvTLbW7gJoVQfg1sh1BmzeNBUs5TsXy0Q7YDgGbfazSZy5LKnU3l", // no password set
			DB:       0,
			PoolSize: 5,
		})

	redisPublisher.client = client

	redisPublisher.pubsub = client.Subscribe(redisPublisher.ctx, xid.New().String())

	go redisPublisher.handleSubscribeMessages()

	return redisPublisher
}
