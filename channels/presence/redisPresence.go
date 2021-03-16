// This package holds the presence interface implementations
package presence

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lisomatrix/channels/channels/core"

	"github.com/go-redis/redis/v8"
)

// RedisPresence - Redis implementation of PresenceHandler
type RedisPresence struct {
	client *redis.Client
	ctx    context.Context
}

// GetChannelClientsPresence - Get channel current presence data
func (presence *RedisPresence) GetChannelClientsPresence(appID string, channelID string) map[string]int64 {
	cmd := presence.client.HGetAll(presence.ctx, appID+":channel:"+channelID+":presence")

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to get channel presences %v\n", cmd.Err())
		return nil
	}

	devicesPresences, err := cmd.Result()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to get channel presences result %v\n", cmd.Err())
		return nil
	}

	lastClientPresences := make(map[string]int64)

	for clientAndDevice, timestampStr := range devicesPresences {
		parts := strings.Split(clientAndDevice, ":")

		clientID := parts[0]

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to convert timestamp %v\n", err)
		}

		if val, ok := lastClientPresences[clientID]; ok {

			if timestamp > val {
				lastClientPresences[clientID] = timestamp
			}

		} else {
			lastClientPresences[clientID] = timestamp
		}
	}

	return lastClientPresences
}

// IsClientDeviceConnectToChannel - Check if device is connected to channel, this will be mostly used to prevent online status constantly changing
func (presence *RedisPresence) IsClientDeviceConnectToChannel(appID string, channelID string, clientID string, deviceID string) bool {
	// HEXISTS channel:{channelID}:presence {clientID}:{deviceID}
	cmd := presence.client.HExists(presence.ctx, appID+":channel:"+channelID+":presence", clientID+":"+deviceID)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to check device is online in channel %v\n", cmd.Err())
		return false
	}

	result, err := cmd.Result()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to check device is online in channel result %v\n", cmd.Err())
		return false
	}

	return result
}

// AddOnlineChannelDevice - Add device to channel
func (presence *RedisPresence) AddOnlineChannelDevice(appID string, channelID string, clientID string, deviceID string) {
	// HSET channel:{channelID}:presence {clientID}:{deviceID} {timestamp}
	cmd := presence.client.HSet(presence.ctx, appID+":channel:"+channelID+":presence", clientID+":"+deviceID, time.Now().Unix())

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to add device to online channel %v\n", cmd.Err())
		return
	}
}

// RemoveOnlineChannelDevice - Remove device from channel
func (presence *RedisPresence) RemoveOnlineChannelDevice(appID string, channelID string, clientID string, deviceID string) {
	// HDEL channel:{channelID}:presence {clientID}:{deviceID}
	cmd := presence.client.HDel(presence.ctx, appID+":channel:"+channelID+":presence", clientID+":"+deviceID)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to remove device from online channel %v\n", cmd.Err())
		return
	}
}

// GetChannelAmountOfClientDevices - Get how many client devices are subscribed to this channel
func (presence *RedisPresence) GetChannelAmountOfClientDevices(appID string, channelID string, clientID string) int64 {
	// HSCAN channel:{channelID}:presence 0 match {clientID}:*
	cmd := presence.client.HScan(presence.ctx, appID+":channel:"+channelID+":presence", 0, clientID+":*", 0)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to get client devices connected to channel %v\n", cmd.Err())
		return 0
	}

	result, _, err := cmd.Result()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to get client devices connected to channel result %v\n", err)
		return 0
	}

	return int64(len(result) - 1)
}

// SetDeviceOnline - Set device online
func (presence *RedisPresence) SetDeviceOnline(clientID string, deviceID string) {
	cmd := presence.client.SAdd(presence.ctx, formatKeyOnline(clientID), deviceID)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to set client device online status %v\n", cmd.Err())
		return
	}
}

// SetDeviceOffline - Set device offline
func (presence *RedisPresence) SetDeviceOffline(clientID string, deviceID string) {
	cmd := presence.client.SRem(presence.ctx, formatKeyOnline(clientID), deviceID)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to set client device online status %v\n", cmd.Err())
		return
	}
}

// UpdateDeviceTimestamp - Update client device connected status
func (presence *RedisPresence) UpdateDeviceTimestamp(clientID string, deviceID string) {
	cmd := presence.client.HSet(presence.ctx, formatKeyPresence(clientID), deviceID, time.Now().Unix())

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to update client device online status %v\n", cmd.Err())
		return
	}
}

// GetClientDevicesPresences -  Get all connected devices with their last online timstamp
func (presence *RedisPresence) GetClientDevicesPresences(clientID string) ([]*core.LastDevicePresence, error) {
	cmd := presence.client.HGetAll(presence.ctx, formatKeyPresence(clientID))

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to retrieve client devices with timestamps %v\n", cmd.Err())
		return nil, cmd.Err()
	}

	result, err := cmd.Result()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to retrieve client devices with timestamps result %v\n", cmd.Err())
		return nil, err
	}

	devices := make([]*core.LastDevicePresence, 0, len(result))

	for key, value := range result {

		timestamp, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to parse timestamp %v\n", cmd.Err())
			return nil, cmd.Err()
		}

		devices = append(devices, &core.LastDevicePresence{
			ClientID:  clientID,
			DeviceID:  key,
			Timestamp: timestamp,
		})
	}

	return devices, nil
}

// GetClientOnlineDevices - Get all connected client devices
func (presence *RedisPresence) GetClientOnlineDevices(clientID string) ([]string, error) {
	cmd := presence.client.SMembers(presence.ctx, formatKeyOnline(clientID))

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to retrieve client devices %v\n", cmd.Err())
		return nil, cmd.Err()
	}

	result, err := cmd.Result()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to retrieve client devices result %v\n", cmd.Err())
	}

	return result, nil
}

// RemoveDevice - Remove client device from connected status
func (presence *RedisPresence) RemoveDevice(clientID string, deviceID string) {
	cmd := presence.client.HDel(presence.ctx, formatKeyPresence(clientID), deviceID)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to remove client device %v\n", cmd.Err())
		return
	}
}

// AddDevice - Add client device to connected status
func (presence *RedisPresence) AddDevice(clientID string, deviceID string) {
	cmd := presence.client.HSet(presence.ctx, formatKeyPresence(clientID), deviceID, time.Now().Unix())

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to add client device %v\n", cmd.Err())
		return
	}
}

// IsOnline - Check if client is online by checking connected devices
func (presence *RedisPresence) IsOnline(clientID string) bool {
	cmd := presence.client.HLen(presence.ctx, formatKeyPresence(clientID))

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to verify client online presence %v\n", cmd.Err())
		return false
	}

	result, err := cmd.Result()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "RedisPresence: failed to retrieve client online presence result %v\n", cmd.Err())
	}

	return result > 0
}

func formatKeyChannelOnline(channelID string, clientID string) string {
	return "channel:" + channelID + ":client:" + clientID
}

func formatKeyOnline(clientID string) string {
	return "client" + ":" + clientID + ":online"
}

func formatKeyPresence(clientID string) string {
	return "client" + ":" + clientID + ":presence"
}

// NewRedisPresence - Create new instance of RedisPresence
func NewRedisPresence() *RedisPresence {
	redisPresence := new(RedisPresence)

	redisPresence.ctx = context.Background()

	redisPresence.client = redis.NewClient(
		&redis.Options{
			Addr:     "127.0.0.1:6379",
			DB:       0,
			PoolSize: 5,
		})

	return redisPresence
}
