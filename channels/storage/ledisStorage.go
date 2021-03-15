package storage

import (
	"encoding/json"
	"github.com/lisomatrix/channels/channels/core"
	"log"

	lediscfg "github.com/ledisdb/ledisdb/config"
	ledis "github.com/ledisdb/ledisdb/ledis"
)

// LedisStorage - Storage implementation with ledis
type LedisStorage struct {
	db *ledis.DB
}

// GetClients - Get all stored clients
func (storage *LedisStorage) GetClients() []*core.Client {
	clients := make([]*core.Client, 0)

	return clients
}

// DeleteClient - Remove client info
func (storage *LedisStorage) DeleteClient(clientID string) {
	_, _ = storage.db.Del([]byte(clientID))
}

// StoreClient - Add new client info
func (storage *LedisStorage) StoreClient(clientID string, client *core.Client) {
	data, err := json.Marshal(client)

	if err != nil {
		log.Println(err)
		return
	}

	_ = storage.db.Set([]byte(clientID), data)
}

// LoadClient - Load client info
func (storage *LedisStorage) LoadClient(clientID string) *core.Client {
	data, err := storage.db.Get([]byte(clientID))

	if err != nil {
		return nil
	}

	var client core.Client

	err = json.Unmarshal(data, &client)

	if err != nil {
		log.Println(err)
		return nil
	}

	return &client
}

// StoreChannel - Add new channel info
func (storage *LedisStorage) StoreChannel(id string, chann *core.Channel) {
	data, err := json.Marshal(chann)

	if err != nil {
		log.Println(err)
		return
	}

	_ = storage.db.Set([]byte(id), data)
}

// GetChannel - Get channel info
func (storage *LedisStorage) GetChannel(id string) *core.Channel {
	data, err := storage.db.Get([]byte(id))

	if err != nil || data == nil {
		return nil
	}

	var chann core.Channel

	err = json.Unmarshal(data, &chann)

	if err != nil {
		log.Println(err)
		return nil
	}

	return &chann
}

// StoreMessage - Store new channel message
func (storage *LedisStorage) StoreMessage(channelID string, event *core.ChannelEvent) {

	data, err := json.Marshal(event)

	if err != nil {
		log.Println(err)
		return
	}

	_, _ = storage.db.LPush([]byte(channelID), data)
}

/*
// GetMessagesSince - Get messages in channel since given timestamp
func (storage *LedisStorage) GetMessagesSince(channelID string, timestamp time.Time) []*core.ChannelEvent {
	events := make([]*core.ChannelEvent, 0)

	dataArr, err := storage.db.LRange([]byte(channelID), 0, -1)

	if err != nil {
		return events
	}

	for _, data := range dataArr {
		var event core.ChannelEvent

		err := json.Unmarshal(data, &event)

		if err != nil {
			log.Println(err)
			continue
		}

		if event.Timestamp.After(timestamp) {
			events = append(events, &event)
		}
	}

	return events
}*/

// NewLedisStorage - Create Storage implementation with ledis
func NewLedisStorage() *LedisStorage {
	cfg := lediscfg.NewConfigDefault()

	l, _ := ledis.Open(cfg)
	db, _ := l.Select(0)

	return &LedisStorage{db: db}
}
