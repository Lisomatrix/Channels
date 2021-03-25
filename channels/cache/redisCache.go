// This package holds the implementations for the project cache interface
package cache

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lisomatrix/channels/channels/core"

	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

// RedisCacheStorage - Cache implementation in Redis
type RedisCacheStorage struct {
	db  *redis.Client
	ctx context.Context
}

// GetChannelEvents - Get given cached events from the channel queue
func (cache *RedisCacheStorage) GetChannelEvents(channelID string, appID string, amount int64) []*core.ChannelEvent {
	cmd := cache.db.LRange(cache.ctx, "app:" + appID + ":channel:" + channelID + ":events", 0, amount)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed get cached events %v\n", cmd.Err())
		return nil
	}

	dData, err := cmd.Result()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed get LRANGE result %v\n", cmd.Err())
		return nil
	}

	events := make([]*core.ChannelEvent, 0, len(dData))

	for _, data := range dData {
		var cachedEvent CachedChannelEvent

		err = proto.Unmarshal([]byte(data), &cachedEvent)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed umarshal cached event %v\n", cmd.Err())
			return nil
		}

		events = append(events, &core.ChannelEvent{
			SenderID:  cachedEvent.SenderID,
			Payload:   cachedEvent.Payload,
			EventType: cachedEvent.EventType,
			ChannelID: channelID,
			Timestamp: cachedEvent.Timestamp,
		})
	}

	return events
}

// GetChannelEventsSize - Get how much events are stored in cache
func (cache *RedisCacheStorage) GetChannelEventsSize(channelID string, appID string) uint64 {
	key := "app:" + appID + ":channel:" + channelID + ":events"
	cmd := cache.db.LLen(cache.ctx, key)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed get cached event queue size %v\n", cmd.Err())
		return 0
	}

	result, err := cmd.Result()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed get LRANGE result %v\n", cmd.Err())
		return 0
	}

	return uint64(result)
}

// GetOldestChannelEvent - Get oldest event that is stored in cache
func (cache *RedisCacheStorage) GetOldestChannelEvent(channelID string, appID string) *core.ChannelEvent {
	cmd := cache.db.LRange(cache.ctx, "app:" + appID + ":channel:" + channelID + ":events", -1, -1)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed get oldest channel event %v\n", cmd.Err())
		return nil
	}

	results, err := cmd.Result()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed get LRANGE result %v\n", cmd.Err())
		return nil
	}

	if len(results) != 1 {
		return nil
	}

	var cachedEvent CachedChannelEvent

	err = proto.Unmarshal([]byte(results[0]), &cachedEvent)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed umarshal cached event %v\n", cmd.Err())
		return nil
	}

	return &core.ChannelEvent{
		SenderID:  cachedEvent.SenderID,
		Payload:   cachedEvent.Payload,
		EventType: cachedEvent.EventType,
		ChannelID: channelID,
		Timestamp: cachedEvent.Timestamp,
	}
}

// StoreChannelEvent - Store channel event on the beginning of the channel queue
func (cache *RedisCacheStorage) StoreChannelEvent(channelID string, appID string, event *core.ChannelEvent) {

	cachedEvent := CachedChannelEvent{
		SenderID:  event.SenderID,
		Payload:   event.Payload,
		Timestamp: event.Timestamp,
		EventType: event.EventType,
	}

	eventData, err := proto.Marshal(&cachedEvent)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to marshal cached event %v\n", err)
		return
	}

	key := "app:" + appID + ":channel:" + channelID + ":events"

	// Push new event and update expire period
	pipeliner := cache.db.Pipeline()
	cmd := pipeliner.LPush(cache.ctx, key, eventData)
	pipeliner.Expire(cache.ctx, key, 24*time.Hour)
	_, err = pipeliner.Exec(cache.ctx)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to pipeline LPUSH and EXPIRE %v\n", err)
		return
	}

	amount, err := cmd.Result()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed get LPUSH result %v\n", cmd.Err())
		return
	}

	if amount > core.CacheQueueSize {
		cache.db.LTrim(cache.ctx, key, 0, core.CacheQueueSize)
	}

}

// CheckDeviceExistence - Check if device exists in cache
func (cache *RedisCacheStorage) CheckDeviceExistence(clientID string, id string) bool {
	cmd := cache.db.HExists(cache.ctx, clientID+":device", id)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to check device existence %v\n", cmd.Err())
		return false
	}

	return cmd.Val()
}

// RemoveDevice - Remove device from user list
func (cache *RedisCacheStorage) RemoveDevice(clientID string, id string) {
	cmd := cache.db.HDel(cache.ctx, clientID+":device", id)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to remove device %v\n", cmd.Err())
	}
}

// AddDevice - Add device to user list
func (cache *RedisCacheStorage) AddDevice(clientID string, device *core.Device) {
	cmd := cache.db.HSet(cache.ctx, clientID+":device", device.ID, device.Token)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to add device %v\n", cmd.Err())
	}
}

// GetClientDevices - Get all client devices
func (cache *RedisCacheStorage) GetClientDevices(clientID string) []*core.Device {
	cmd := cache.db.HGetAll(cache.ctx, clientID+":device")

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: to get devices %v\n", cmd.Err())
		return nil
	}

	data, err := cmd.Result()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to get devices data %v\n", err)
		return nil
	}

	devices := make([]*core.Device, 0, len(data))

	for deviceID, token := range data {

		devices = append(devices, &core.Device{
			ID:       deviceID,
			Token:    token,
			ClientID: clientID,
		})
	}

	return devices
}

// RemoveClient - Remove client from cache
func (cache *RedisCacheStorage) RemoveClient(appID string, clientID string) {
	cmd := cache.db.Del(cache.ctx, appID+":client:"+clientID)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to remove client %v\n", cmd.Err())
	}
}

// RemoveClientChannels - Remove client channels from cache
func (cache *RedisCacheStorage) RemoveClientChannels(clientID string) {
	cmd := cache.db.Del(cache.ctx, "client:"+clientID+":channels")

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to remove client channels %v\n", cmd.Err())
	}
}

// RemoveApp - Remove app from cache
func (cache *RedisCacheStorage) RemoveApp(appID string) {
	cmd := cache.db.Del(cache.ctx, "app:"+appID)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to remove app %v\n", cmd.Err())
	}
}

// RemoveChannel - Remove channel from cache
func (cache *RedisCacheStorage) RemoveChannel(appID string, channelID string) {
	cmd := cache.db.Del(cache.ctx, appID+":channel:"+channelID)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to remove channel %v\n", cmd.Err())
	}
}

// StoreClient - Cache client
func (cache *RedisCacheStorage) StoreClient(appID string, clientID string, client *core.Client) {

	// Lunch a goroutine to reduce latency
	go func() {
		err := cache.db.HMSet(cache.ctx, appID+":client:"+clientID, "username", client.Username, "extra", client.Extra)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to store client %v\n", err)
			return
		}
	}()
}

// CheckClientExistence - Check if there is a client in cache
func (cache *RedisCacheStorage) CheckClientExistence(appID string, clientID string) bool {
	cmd := cache.db.Exists(cache.ctx, appID+":client:"+clientID)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to check client existence %v\n", cmd.Err())
		return false
	}

	return cmd.Val() > 0
}

// GetClient - Attempt to get client from cache
func (cache *RedisCacheStorage) GetClient(appID string, clientID string) *core.Client {
	cmd := cache.db.HMGet(cache.ctx, appID+":client:"+clientID, "username", "extra")

	dData, err := cmd.Result()

	if cmd.Err() != nil || len(dData) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to retrieve cached client %v\n", err)
		return nil
	}

	var cachedClient CachedClient

	err = cmd.Scan(&cachedClient)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to umarshal cached client %v\n", err)
		return nil
	}

	// If username is empty it wasn't in cache
	if cachedClient.Username == "" {
		return nil
	}

	return &core.Client{
		ID:       clientID,
		AppID:    appID,
		Username: cachedClient.Username,
		Extra:    cachedClient.Extra,
	}
}

// StoreApp - Set app in cache
func (cache *RedisCacheStorage) StoreApp(appID string, name string) {
	cmd := cache.db.Set(cache.ctx, "app:"+appID, name, 0)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to store cached app %v\n", cmd.Err())
	}
}

// GetApp - Get app from cache
func (cache *RedisCacheStorage) GetApp(appID string) *core.App {
	cmd := cache.db.Get(cache.ctx, "app:"+appID)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to retrieve cached app %v\n", cmd.Err())
		return nil
	}

	data, err := cmd.Result()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to parse cached app %v\n", err)
		return nil
	}

	return &core.App{
		AppID: appID,
		Name:  string(data),
	}
}

// StoreChannel - Store channel in cache
func (cache *RedisCacheStorage) StoreChannel(appID string, channelID string, channel *core.Channel) {

	// Lunch a goroutine to reduce latency
	go func() {
		cachedChannel := CachedChannel{
			Name:       channel.Name,
			CreatedAt:  channel.CreatedAt,
			IsClosed:   channel.IsClosed,
			Extra:      channel.Extra,
			Persistent: channel.Persistent,
			Private:    channel.Private,
			Presence:   channel.Presence,
		}

		data, err := proto.Marshal(&cachedChannel)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to marshal cached channel %v\n", err)
			return
		}

		cmd := cache.db.Set(cache.ctx, appID+":channel:"+channelID, data, 0)

		if cmd.Err() != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to store cached channel %v\n", err)
			return
		}
	}()
}

// GetChannel - Get channel from cache
func (cache *RedisCacheStorage) GetChannel(appID string, channelID string) *core.Channel {
	cmd := cache.db.Get(cache.ctx, appID+":channel:"+channelID)

	// In case it has nil then there is ni need to log this data, just means it didn't find the key
	if cmd.Err() != nil && strings.Contains(cmd.Err().Error(), "nil") {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to retrieve execute command %v\n", cmd.Err())
		return nil
	}

	data, err := cmd.Bytes()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to retrieve cached channel %v\n", err)
		return nil
	}

	var cachedChannel CachedChannel

	err = proto.Unmarshal(data, &cachedChannel)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to umarshal cached channel %v\n", err)
		return nil
	}

	return &core.Channel{
		ID:         channelID,
		AppID:      appID,
		Name:       cachedChannel.Name,
		Extra:      cachedChannel.Extra,
		CreatedAt:  cachedChannel.CreatedAt,
		IsClosed:   cachedChannel.IsClosed,
		Persistent: cachedChannel.Persistent,
		Private:    cachedChannel.Private,
		Presence:   cachedChannel.Presence,
	}
}

// CheckChannelExistence - Check if there is channel in cache
func (cache *RedisCacheStorage) CheckChannelExistence(appID string, channelID string) bool {
	cmd := cache.db.Exists(cache.ctx, appID+":channel:"+channelID)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to check cached channel existence %v\n", cmd.Err())
		return false
	}

	amount, err := cmd.Uint64()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to verify channel existence %v\n", cmd.Err())
		return false
	}

	return amount > 0
}

// AddClientChannels - Store list of channels client can access in cache
func (cache *RedisCacheStorage) AddClientChannels(clientID string, channelIDs []string) {

	channelsBin := make([][]byte, 0, len(channelIDs))

	for _, c := range channelIDs {
		channelsBin = append(channelsBin, []byte(c))
	}

	cmd := cache.db.SAdd(cache.ctx, "client:"+clientID+":channels", channelsBin)

	if cmd.Err() != nil {
		fmt.Fprintf(os.Stderr, "Redis Cache: failed to add multiple client channels %v\n", cmd.Err())
	}

}

// GetClientChannels - Get channels client can access from cache
func (cache *RedisCacheStorage) GetClientChannels(clientID string) ([]string, bool) {
	cmd := cache.db.SMembers(cache.ctx, "client:"+clientID+":channels")

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to executre command %v\n", cmd.Err())
		return nil, false
	}

	dData, err := cmd.Result()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to retrieve client channels %v\n", err)
		return nil, false
	}

	if len(dData) == 0 {
		return nil, false
	}

	channels := make([]string, 0, len(dData))

	for _, data := range dData {
		channels = append(channels, data)
	}

	return channels, true
}

// AddClientChannel - Add a new channel to client channels cache
func (cache *RedisCacheStorage) AddClientChannel(clientID string, channelID string) {
	cmd := cache.db.SAdd(cache.ctx, "client:"+clientID+":channels", channelID)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to add single client channel %v\n", cmd.Err())
	}
}

// RemoveClientChannel - Remove a channel from client channels cache
func (cache *RedisCacheStorage) RemoveClientChannel(clientID string, channelID string) {
	cmd := cache.db.SRem(cache.ctx, "client:"+clientID+":channels", channelID)

	if cmd.Err() != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Redis Cache: failed to store client channels %v\n", cmd.Err())
	}
}


// NewRedisCacheStorage - Create a new Redis cache instance
func NewRedisCacheStorage() *RedisCacheStorage {

	db := redis.NewClient(
		&redis.Options{
			Addr: "127.0.0.1:6379",
			// Password: "CulP3gnpgSAxFlbjO/JrNCR/uTKFKvTLbW7gJoVQfg1sh1BmzeNBUs5TsXy0Q7YDgGbfazSZy5LKnU3l", // no password set
			DB:       0,
			PoolSize: 5,
		})

	return &RedisCacheStorage{db: db, ctx: context.Background()}
}
