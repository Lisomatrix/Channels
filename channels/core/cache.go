package core

// CacheQueueSize - How much the cache of a channel queue can grow, the bigger the less database request we make
var CacheQueueSize int64 = 50

// CacheStorage - Cache for avoiding fetching data from database
type CacheStorage interface {
	// Device
	CheckDeviceExistence(clientID string, id string) bool
	GetClientDevices(clientID string) []*Device
	RemoveDevice(clientID string, id string)
	AddDevice(clientID string, device *Device)
	// Client
	StoreClient(appID string, clientID string, client *Client)
	CheckClientExistence(appID string, clientID string) bool
	GetClient(appID string, clientID string) *Client
	RemoveClient(appID string, clientID string)
	// App
	StoreApp(appID string, name string)
	GetApp(appID string) *App
	RemoveApp(appID string)
	// Channel
	StoreChannel(appID string, channelID string, channel *Channel)
	GetChannel(appID string, channelID string) *Channel
	CheckChannelExistence(appID string, channelID string) bool
	RemoveChannel(appID string, channelID string)
	// Client -> Channel
	RemoveClientChannels(clientID string)
	AddClientChannels(clientID string, channelIDs []string)
	AddClientChannel(clientID string, channelID string)
	GetClientChannels(clientID string) ([]string, bool)
	RemoveClientChannel(clientID string, channelID string)
	// Channel Event
	StoreChannelEvent(channelID string, event *ChannelEvent)
	GetOldestChannelEvent(channelID string) *ChannelEvent
	GetChannelEventsSize(channelID string) uint64
	GetChannelEvents(channelID string, amount int64) []*ChannelEvent
}

// REDIS APP

// Store app
// hset app:id name
// Get App
// hget app:id

// REDIS CHANNEL

// Store channel
// set appID:channel:ID json/protobuf
// Get channel
// get appID:channel:ID
// Channel exists
// exists appID:channel:ID

// REDIS CLIENT

// Store client - storing appID might not be useful
// hmset client:id username :username appID :appID extra :extra
// OR
// set appID:client:id json/protobuf

// Get client
// hmget appID:client:id username appID extra
// OR
// get appID:client:id

// Check client existence
// exists appID:client:id

// Store client channels
// sadd client:id:channels :channelID1 :channelID2 :channelID3
// Remove client channel
// srem client:id:channels :channelID1
// Get client channels
// smembers client:id:channels
