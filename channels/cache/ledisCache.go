package cache

import (
	"fmt"
	"github.com/lisomatrix/channels/channels/core"
	"os"
	"time"

	lediscfg "github.com/ledisdb/ledisdb/config"
	"github.com/ledisdb/ledisdb/ledis"
	"google.golang.org/protobuf/proto"
)

// LedisCacheStorage - Cache implementation in Ledis
type LedisCacheStorage struct {
	db *ledis.DB
}

// GetChannelEvents - Get given cached events from the channel queue
func (cache *LedisCacheStorage) GetChannelEvents(channelID string, appID string, amount int64) []*core.ChannelEvent {
	dData, err := cache.db.LRange([]byte("app:" + appID + "channel:" + channelID + ":events"), 0, int32(amount))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed get LRANGE result %v\n", err)
		return nil
	}

	if dData == nil {
		return nil
	}

	events := make([]*core.ChannelEvent, 0, len(dData))

	for _, data := range dData {
		var cachedEvent CachedChannelEvent

		err = proto.Unmarshal([]byte(data), &cachedEvent)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed umarshal cached event %v\n", err)
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
func (cache *LedisCacheStorage) GetChannelEventsSize(channelID string, appID string) uint64 {
	key := []byte("app:" + appID + "channel:" + channelID + ":events")

	result, err := cache.db.LLen(key)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed get LRANGE result %v\n", err)
		return 0
	}

	return uint64(result)
}

// GetOldestChannelEvent - Get oldest event that is stored in cache
func (cache *LedisCacheStorage) GetOldestChannelEvent(channelID string, appID string) *core.ChannelEvent {
	results, err := cache.db.LRange([]byte("app:" + appID + ":channel:" + channelID + ":events"), -1, -1)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed get LRANGE result %v\n", err)
		return nil
	}

	if len(results) != 1 {
		return nil
	}

	var cachedEvent CachedChannelEvent

	err = proto.Unmarshal([]byte(results[0]), &cachedEvent)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed umarshal cached event %v\n", err)
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
func (cache *LedisCacheStorage) StoreChannelEvent(channelID string, appID string, event *core.ChannelEvent) {

	cachedEvent := CachedChannelEvent{
		SenderID:  event.SenderID,
		Payload:   event.Payload,
		Timestamp: event.Timestamp,
		EventType: event.EventType,
	}

	eventData, err := proto.Marshal(&cachedEvent)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to marshal cached event %v\n", err)
		return
	}

	key := []byte("app:" + appID + "channel:" + channelID + ":events")

	// Push new event and update expire period
	amount, err := cache.db.LPush(key, eventData)

	_, _ = cache.db.Expire(key, int64((24 * time.Hour).Seconds()))


	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed get LPUSH result %v\n", err)
		return
	}

	if amount > core.CacheQueueSize {
		_ = cache.db.LTrim(key, 0, core.CacheQueueSize)
	}

}

// CheckDeviceExistence - Check if device exists in cache
func (cache *LedisCacheStorage) CheckDeviceExistence(clientID string, id string) bool {
	amount, err := cache.db.HGet([]byte(clientID+":device"), []byte(id))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to check device existence %v\n", err)
		return false
	}

	return amount != nil
}

// RemoveDevice - Remove device from user list
func (cache *LedisCacheStorage) RemoveDevice(clientID string, id string) {
	_, err := cache.db.HDel([]byte(clientID+":device"), []byte(id))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to remove device %v\n", err)
	}
}

// AddDevice - Add device to user list
func (cache *LedisCacheStorage) AddDevice(clientID string, device *core.Device) {
	_, err := cache.db.HSet([]byte(clientID+":device"), []byte(device.ID), []byte(device.Token))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to add device %v\n", err)
	}
}

// GetClientDevices - Get all client devices
func (cache *LedisCacheStorage) GetClientDevices(clientID string) []*core.Device {
	data, err  := cache.db.HGetAll([]byte(clientID+":device"))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to get devices data %v\n", err)
		return nil
	}

	devices := make([]*core.Device, 0, len(data))

	for _, pair := range data {

		devices = append(devices, &core.Device{
			ID:       string(pair.Field),
			Token:    string(pair.Value),
			ClientID: clientID,
		})
	}

	return devices
}

// RemoveClient - Remove client from cache
func (cache *LedisCacheStorage) RemoveClient(appID string, clientID string) {
	_, err := cache.db.Del([]byte(appID+":client:"+clientID))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to remove client %v\n", err)
	}
}

// RemoveClientChannels - Remove client channels from cache
func (cache *LedisCacheStorage) RemoveClientChannels(clientID string) {
	_, err := cache.db.Del([]byte("client:"+clientID+":channels"))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to remove client channels %v\n", err)
	}
}

// RemoveApp - Remove app from cache
func (cache *LedisCacheStorage) RemoveApp(appID string) {
	_, err := cache.db.Del([]byte("app:"+appID))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to remove app %v\n", err)
	}
}

// RemoveChannel - Remove channel from cache
func (cache *LedisCacheStorage) RemoveChannel(appID string, channelID string) {
	_, err := cache.db.Del([]byte(appID+":channel:"+channelID))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to remove channel %v\n", err)
	}
}

// StoreClient - Cache client
func (cache *LedisCacheStorage) StoreClient(appID string, clientID string, client *core.Client) {

	// Lunch a goroutine to reduce latency
	go func() {

		err := cache.db.HMset(
				[]byte(appID+":client:"+clientID),
				ledis.FVPair{Value: []byte(client.Username), Field: []byte("username")},
				ledis.FVPair{Value: []byte(client.Extra), Field: []byte("extra")},
			)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to store client %v\n", err)
			return
		}
	}()
}

// CheckClientExistence - Check if there is a client in cache
func (cache *LedisCacheStorage) CheckClientExistence(appID string, clientID string) bool {
	amount, err := cache.db.Exists([]byte(appID+":client:"+clientID))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to check client existence %v\n", err)
		return false
	}

	return amount > 0
}

// GetClient - Attempt to get client from cache
func (cache *LedisCacheStorage) GetClient(appID string, clientID string) *core.Client {
	dData, err := cache.db.HMget([]byte(appID+":client:"+clientID), []byte("username"), []byte("extra"))

	if err != nil || len(dData) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to retrieve cached client %v\n", err)
		return nil
	}

	if dData == nil {
		return nil
	}

	var cachedClient CachedClient

	for i, data := range dData {

		if i == 1 {
			cachedClient.Username = string(data)
		} else if i == 3 {
			cachedClient.Extra = string(data)
		}
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
func (cache *LedisCacheStorage) StoreApp(appID string, name string) {
	err := cache.db.Set([]byte("app:"+appID), []byte(name))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to store cached app %v\n", err)
	}
}

// GetApp - Get app from cache
func (cache *LedisCacheStorage) GetApp(appID string) *core.App {
	data, err := cache.db.Get([]byte("app:"+appID))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to parse cached app %v\n", err)
		return nil
	}

	if data == nil {
		return nil
	}

	return &core.App{
		AppID: appID,
		Name:  string(data),
	}
}

// StoreChannel - Store channel in cache
func (cache *LedisCacheStorage) StoreChannel(appID string, channelID string, channel *core.Channel) {

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
			Push: 		channel.Push,
		}

		data, err := proto.Marshal(&cachedChannel)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to marshal cached channel %v\n", err)
			return
		}

		err = cache.db.Set([]byte(appID+":channel:"+channelID), data)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to store cached channel %v\n", err)
			return
		}
	}()
}

// GetChannel - Get channel from cache
func (cache *LedisCacheStorage) GetChannel(appID string, channelID string) *core.Channel {
	data, err := cache.db.Get([]byte(appID+":channel:"+channelID))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to retrieve cached channel %v\n", err)
		return nil
	}

	if data == nil {
		return nil
	}

	var cachedChannel CachedChannel

	err = proto.Unmarshal(data, &cachedChannel)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to umarshal cached channel %v\n", err)
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
func (cache *LedisCacheStorage) CheckChannelExistence(appID string, channelID string) bool {
	amount, err := cache.db.Exists([]byte(appID+":channel:"+channelID))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to verify channel existence %v\n", err)
		return false
	}

	return amount > 0
}

// AddClientChannels - Store list of channels client can access in cache
func (cache *LedisCacheStorage) AddClientChannels(clientID string, channelIDs []string) {

	channelsBin := make([][]byte, 0, len(channelIDs))

	for _, c := range channelIDs {
		channelsBin = append(channelsBin, []byte(c))
	}

	_, err := cache.db.SAdd([]byte("client:"+clientID+":channels"), channelsBin...)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to add multiple client channels %v\n", err)
	}

}

// GetClientChannels - Get channels client can access from cache
func (cache *LedisCacheStorage) GetClientChannels(clientID string) ([]string, bool) {
	dData, err := cache.db.SMembers([]byte("client:"+clientID+":channels"))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to retrieve client channels %v\n", err)
		return nil, false
	}

	if len(dData) == 0 {
		return nil, false
	}

	channels := make([]string, 0, len(dData))

	for _, data := range dData {
		channels = append(channels, string(data))
	}

	return channels, true
}

// AddClientChannel - Add a new channel to client channels cache
func (cache *LedisCacheStorage) AddClientChannel(clientID string, channelID string) {
	_, err := cache.db.SAdd([]byte("client:"+clientID+":channels"), []byte(channelID))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to add single client channel %v\n", err)
	}
}

// RemoveClientChannel - Remove a channel from client channels cache
func (cache *LedisCacheStorage) RemoveClientChannel(clientID string, channelID string) {
	_, err := cache.db.SRem([]byte("client:"+clientID+":channels"), []byte(channelID))

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ledis Cache: failed to store client channels %v\n", err)
	}
}


func (cache *LedisCacheStorage) GetDB() *ledis.DB {
	return cache.db
}


// NewLedisCacheStorage - Create a new ledis cache instance
func NewLedisCacheStorage() *LedisCacheStorage {

	cfg := lediscfg.NewConfigDefault()

	l, _ := ledis.Open(cfg)
	db, _ := l.Select(0)

	// We need to delete all cache, or we could get old values
	_, _ = db.FlushAll()

	return &LedisCacheStorage{db: db}
}
