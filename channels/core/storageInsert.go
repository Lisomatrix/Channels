package core

import (
	log "github.com/sirupsen/logrus"
)

// InsertItem - Insert Item queued
type InsertItem struct {
	Event *ChannelEvent
	AppID string
}

type StorageInsert interface {
	StoreEvent(appID string, event *ChannelEvent)
	Start(channelRepository ChannelRepository)
}

func NewStorageInsertQueue() *StorageInsertQueue {
	return &StorageInsertQueue{
		insertChannel: make(chan InsertItem, 300),
	}
}

// StorageInsertQueue - Receives all insert requests and send them into the database
type StorageInsertQueue struct {
	insertChannel chan InsertItem
}

func (storage *StorageInsertQueue) StoreEvent(appID string, event *ChannelEvent) {
	storage.insertChannel <- InsertItem{
		AppID: appID,
		Event: event,
	}
}

func (storage *StorageInsertQueue) Start(repo ChannelRepository) {
	for {
		item, isActive := <-storage.insertChannel

		if !isActive {
			return
		}

		if err := repo.AddChannelEvent(item.AppID, item.Event.ChannelID, item.Event); err != nil {
			log.WithFields(log.Fields{
				"Item": item,
			}).Error(err)
		}
	}
}

// func (storage *StorageInsertQueue) Startx(repo ChannelRepository) {
// 	cache := make([]InsertItem, 0, CacheLimit)
// 	tick := time.NewTicker(CacheTimeout)

// 	for {
// 		select {
// 		case m := <-storage.insertChannel:
// 			cache = append(cache, m)

// 			if len(cache) < CacheLimit {
// 				break
// 			}

// 			// Reset the timeout ticker.
// 			// Otherwise we will get too many sends.
// 			tick.Stop()

// 			// Send the cached messages and reset the cache.
// 			if err := repo.AddChannelEvents(cache); err != nil {
// 				log.Printf("error: %v", err)
// 			}
// 			cache = cache[:0]

// 			// Recreate the ticker, so the timeout trigger
// 			// remains consistent.
// 			tick = time.NewTicker(CacheTimeout)
// 		case <-tick.C:
// 			if len(cache) == 0 {
// 				continue
// 			}

// 			if err := repo.AddChannelEvents(cache); err != nil {
// 				log.Printf("error: %v", err)
// 			}
// 			cache = cache[:0]
// 		}
// 	}
// }
