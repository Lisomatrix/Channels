package cache

import (
	"fmt"
	"os"

	"github.com/channelserver/channelserver/core"
	lediscfg "github.com/ledisdb/ledisdb/config"
	ledis "github.com/ledisdb/ledisdb/ledis"
	"google.golang.org/protobuf/proto"
)

// LedisCacheStorage - Cache implementation in Ledis
type LedisCacheStorage struct {
	db *ledis.DB
}

// RemoveClient - Remove client from cache
func (cache *LedisCacheStorage) RemoveClient(appID string, clientID string) {
	_, err := cache.db.Del([]byte(appID + ":client:" + clientID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to remove client %v\n", err)
	}
}

// RemoveClientChannels - Remove client channels from cache
func (cache *LedisCacheStorage) RemoveClientChannels(clientID string) {
	_, err := cache.db.Del([]byte("client:" + clientID + ":channels"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to remove client channels %v\n", err)
	}
}

// RemoveApp - Remove app from cache
func (cache *LedisCacheStorage) RemoveApp(appID string) {
	_, err := cache.db.Del([]byte("app:" + appID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to remove app %v\n", err)
	}
}

// RemoveChannel - Remove channel from cache
func (cache *LedisCacheStorage) RemoveChannel(appID string, channelID string) {
	_, err := cache.db.Del([]byte(appID + ":channel:" + channelID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to remove channel %v\n", err)
	}
}

// StoreClient - Cache client
func (cache *LedisCacheStorage) StoreClient(appID string, clientID string, client *core.Client) {
	// cachedClient := CachedClient{
	// 	Username: client.Username,
	// 	Extra:    client.Extra,
	// }

	// data, err := proto.Marshal(&cachedClient)

	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Ledis Cache: marshal cache %v\n", err)
	// 	return
	// }

	// err = cache.db.Set([]byte(appID+":client:"+clientID), data)

	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Ledis Cache: failed to store client %v\n", err)
	// 	return
	// }

	usernameField := ledis.FVPair{
		Field: []byte("username"),
		Value: []byte(client.Username),
	}

	extraField := ledis.FVPair{
		Field: []byte("extra"),
		Value: []byte(client.Extra),
	}

	err := cache.db.HMset([]byte(appID+":client:"+clientID), usernameField, extraField)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to store client %v\n", err)
		return
	}
}

// CheckClientExistence - Check if there is a client in cache
func (cache *LedisCacheStorage) CheckClientExistence(appID string, clientID string) bool {
	amount, err := cache.db.Exists([]byte(appID + ":client:" + clientID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to check client existence %v\n", err)
		return false
	}

	return amount > 0
}

// GetClient - Attempt to get client from cache
func (cache *LedisCacheStorage) GetClient(appID string, clientID string) *core.Client {
	// data, err := cache.db.Get([]byte(appID + ":client:" + clientID))

	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Ledis Cache: failed to retrieve client %v\n", err)
	// 	return nil
	// }

	// var cachedClient CachedClient

	// err = proto.Unmarshal(data, &cachedClient)

	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Ledis Cache: failed to umarshal cached client %v\n", err)
	// 	return nil
	// }

	dData, err := cache.db.HMget([]byte(appID+":client:"+clientID), []byte("username"), []byte("extra"))

	if err != nil || len(dData) != 2 {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to retrieve cached client %v\n", err)
		return nil
	}

	return &core.Client{
		ID:       clientID,
		AppID:    appID,
		Username: string(dData[0]),
		Extra:    string(dData[1]),
	}

	// return &core.Client{
	// 	ID:       clientID,
	// 	AppID:    appID,
	// 	Username: cachedClient.Username,
	// 	Extra:    cachedClient.Extra,
	// }
}

// StoreApp - Set app in cache
func (cache *LedisCacheStorage) StoreApp(appID string, name string) {
	err := cache.db.Set([]byte("app:"+appID), []byte(name))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to store cached app %v\n", err)
	}
}

// GetApp - Get app from cache
func (cache *LedisCacheStorage) GetApp(appID string) *core.App {
	data, err := cache.db.Get([]byte("app:" + appID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to retrieve cached app %v\n", err)
		return nil
	}

	return &core.App{
		AppID: appID,
		Name:  string(data),
	}
}

// StoreChannel - Store channel in cache
func (cache *LedisCacheStorage) StoreChannel(appID string, channelID string, channel *core.Channel) {

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
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to marshal cached channel %v\n", err)
		return
	}

	err = cache.db.Set([]byte(appID+":channel:"+channelID), data)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to store cached channel %v\n", err)
		return
	}
}

// GetChannel - Get channel from cache
func (cache *LedisCacheStorage) GetChannel(appID string, channelID string) *core.Channel {
	data, err := cache.db.Get([]byte(appID + ":channel:" + channelID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to retrieve cached channel %v\n", err)
		return nil
	}

	var cachedChannel CachedChannel

	err = proto.Unmarshal(data, &cachedChannel)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to umarshal cached channel %v\n", err)
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
	amount, err := cache.db.Exists([]byte(appID + ":channel:" + channelID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to check cached channel existence %v\n", err)
		return false
	}

	return amount > 0
}

// AddClientChannels - Store list of channels client can access in cache
func (cache *LedisCacheStorage) AddClientChannels(clientID string, channelIDs []string) {

	channelsBin := make([][]byte, 0)

	for _, c := range channelIDs {
		channelsBin = append(channelsBin, []byte(c))
	}

	_, err := cache.db.SAdd([]byte("client:" + clientID + ":channels"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to add multiple client channels %v\n", err)
	}

}

// GetClientChannels - Get channels client can access from cache
func (cache *LedisCacheStorage) GetClientChannels(clientID string) ([]string, bool) {
	dData, err := cache.db.SMembers([]byte("client:" + clientID + ":channels"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to retrieve client channels %v\n", err)
		return nil, false
	}

	if len(dData) == 0 {
		return nil, false
	}

	channels := make([]string, 0)

	for _, data := range dData {
		channels = append(channels, string(data))
	}

	return channels, true
}

// AddClientChannel - Add a new channel to client channels cache
func (cache *LedisCacheStorage) AddClientChannel(clientID string, channelID string) {
	_, err := cache.db.SAdd([]byte("client:"+clientID+":channels"), []byte(channelID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to add single client channel %v\n", err)
	}
}

// RemoveClientChannel - Remove a channel from client channels cache
func (cache *LedisCacheStorage) RemoveClientChannel(clientID string, channelID string) {
	_, err := cache.db.SRem([]byte("client:"+clientID+":channels"), []byte(channelID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ledis Cache: failed to store client channels %v\n", err)
	}
}

// NewLedisCacheStorage - Create a new ledis cache instance
func NewLedisCacheStorage() *LedisCacheStorage {

	cfg := lediscfg.NewConfigDefault()

	l, _ := ledis.Open(cfg)
	db, _ := l.Select(0)

	return &LedisCacheStorage{db: db}
}
