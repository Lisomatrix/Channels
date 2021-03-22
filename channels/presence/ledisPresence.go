// This package holds the presence interface implementations
package presence

import (
	"fmt"
	lediscfg "github.com/ledisdb/ledisdb/config"
	"github.com/ledisdb/ledisdb/ledis"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lisomatrix/channels/channels/core"
)

// LedisPresence - Ledis implementation of PresenceHandler
type LedisPresence struct {
	client *ledis.DB
}

func (presence *LedisPresence) UpdateClientTimestamp(clientID string) {
	key := []byte("client:" + clientID + ":hb")
	err := presence.client.Set(key, []byte(strconv.FormatInt(time.Now().Unix(), 10)))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to get update client last heartbeat %v\n", err)
	}

	_, _ = presence.client.ExpireAt(key, int64((time.Minute * 1).Seconds()))
}
func (presence *LedisPresence) GetClientTimestamp(clientID string) int64 {
	data, err := presence.client.Get([]byte("client:" + clientID + ":hb"))


	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to get get client last heartbeast %v\n", err)
		return 0
	}

	timestamp, err := strconv.ParseInt(string(data), 10, 64)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to convert timestamp %v\n", err)
		return 0
	}

	return timestamp
}

// GetChannelClientsPresence - Get channel current presence data
func (presence *LedisPresence) GetChannelClientsPresence(appID string, channelID string) map[string]int64 {
	devicesPresences, err := presence.client.HGetAll([]byte(appID+":channel:"+channelID+":presence"))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to get channel presences result %v\n", err)
		return nil
	}

	if devicesPresences == nil {
		return nil
	}

	lastClientPresences := make(map[string]int64)

	for _, pair := range devicesPresences {
		parts := strings.Split(string(pair.Field), ":")

		clientID := parts[0]

		timestamp, err := strconv.ParseInt(string(pair.Value), 10, 64)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to convert timestamp %v\n", err)
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
func (presence *LedisPresence) IsClientDeviceConnectToChannel(appID string, channelID string, clientID string, deviceID string) bool {
	result, err := presence.client.HGet([]byte(appID+":channel:"+channelID+":presence"), []byte(clientID+":"+deviceID))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to check device is online in channel result %v\n", err)
		return false
	}

	return result != nil
}

// AddOnlineChannelDevice - Add device to channel
func (presence *LedisPresence) AddOnlineChannelDevice(appID string, channelID string, clientID string, deviceID string) {
	_, err := presence.client.HSet( []byte(appID+":channel:"+channelID+":presence"), []byte(clientID+":"+deviceID), []byte(strconv.FormatInt(time.Now().Unix(), 10)))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to add device to online channel %v\n", err)
		return
	}
}

// RemoveOnlineChannelDevice - Remove device from channel
func (presence *LedisPresence) RemoveOnlineChannelDevice(appID string, channelID string, clientID string, deviceID string) {
	// HDEL channel:{channelID}:presence {clientID}:{deviceID}
	_, err := presence.client.HDel([]byte(appID+":channel:"+channelID+":presence"), []byte(clientID+":"+deviceID))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to remove device from online channel %v\n", err)
		return
	}
}

// GetChannelAmountOfClientDevices - Get how many client devices are subscribed to this channel
func (presence *LedisPresence) GetChannelAmountOfClientDevices(appID string, channelID string, clientID string) int64 {
	// HSCAN channel:{channelID}:presence 0 match {clientID}:*
	pairs, err := presence.client.HScan([]byte(appID+":channel:"+channelID+":presence"), []byte{}, 0, true, clientID+":*")

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to get client devices connected to channel result %v\n", err)
		return 0
	}

	return int64(len(pairs) - 1)
}

// SetDeviceOnline - Set device online
func (presence *LedisPresence) SetDeviceOnline(clientID string, deviceID string) {
	_, err := presence.client.SAdd([]byte(formatKeyOnline(clientID)), []byte(deviceID))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to set client device online status %v\n", err)
		return
	}
}

// SetDeviceOffline - Set device offline
func (presence *LedisPresence) SetDeviceOffline(clientID string, deviceID string) {
	_, err := presence.client.SRem([]byte(formatKeyOnline(clientID)), []byte(deviceID))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to set client device online status %v\n", err)
		return
	}
}

// UpdateDeviceTimestamp - Update client device connected status
func (presence *LedisPresence) UpdateDeviceTimestamp(clientID string, deviceID string) {
	_, err := presence.client.HSet([]byte(formatKeyPresence(clientID)), []byte(deviceID), []byte(strconv.FormatInt(time.Now().Unix(), 10)))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to update client device online status %v\n", err)
		return
	}
}

// GetClientDevicesPresences -  Get all connected devices with their last online timstamp
func (presence *LedisPresence) GetClientDevicesPresences(clientID string) ([]*core.LastDevicePresence, error) {
	result, err := presence.client.HGetAll([]byte(formatKeyPresence(clientID)))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to retrieve client devices with timestamps result %v\n", err)
		return nil, err
	}

	devices := make([]*core.LastDevicePresence, 0, len(result))

	for _, keyValue := range result {

		timestamp, err := strconv.ParseInt(string(keyValue.Value), 10, 64)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to parse timestamp %v\n", err)
			return nil, err
		}

		devices = append(devices, &core.LastDevicePresence{
			ClientID:  clientID,
			DeviceID:  string(keyValue.Field),
			Timestamp: timestamp,
		})
	}

	return devices, nil
}

// GetClientOnlineDevices - Get all connected client devices
func (presence *LedisPresence) GetClientOnlineDevices(clientID string) ([]string, error) {
	result, err := presence.client.SMembers([]byte(formatKeyOnline(clientID)))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to retrieve client devices result %v\n", err)
	}

	devices := make([]string, len(result))

	for _, id := range result {
		devices = append(devices, string(id))
	}

	return devices, nil
}

// RemoveDevice - Remove client device from connected status
func (presence *LedisPresence) RemoveDevice(clientID string, deviceID string) {
	_, err := presence.client.HDel([]byte(formatKeyPresence(clientID)), []byte(deviceID))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to remove client device %v\n", err)
		return
	}
}

// AddDevice - Add client device to connected status
func (presence *LedisPresence) AddDevice(clientID string, deviceID string) {
	_, err := presence.client.HSet([]byte(formatKeyPresence(clientID)), []byte(deviceID), []byte(strconv.FormatInt(time.Now().Unix(), 10)))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to add client device %v\n", err)
		return
	}
}

// IsOnline - Check if client is online by checking connected devices
func (presence *LedisPresence) IsOnline(clientID string) bool {
	result, err := presence.client.HLen([]byte(formatKeyPresence(clientID)))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LedisPresence: failed to retrieve client online presence result %v\n", err)
	}

	return result > 0
}

func (presence *LedisPresence) GetDB() *ledis.DB {
	return presence.client
}

// NewLedisPresence - Create new instance of LedisPresence
func NewLedisPresence() *LedisPresence {
	cfg := lediscfg.NewConfigDefault()

	l, _ := ledis.Open(cfg)
	db, _ := l.Select(1)


	// We need to delete all cache, or we could get old values
	_, _ = db.FlushAll()

	return &LedisPresence{client: db}
}

// NewLedisPresence - Create new instance of LedisPresence with given db
func NewLedisPresenceWithDB(db *ledis.DB) *LedisPresence {
	return &LedisPresence{client: db}
}
