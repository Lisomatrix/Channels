package core

import (
	"log"
	"time"

	"github.com/rs/xid"
)

// Engine - Holds application components
type Engine struct {
	serverID        string
	HubsHandler     HubsHandler
	databaseStorage DatabaseStorage
	cacheStorage    CacheStorage
	insertQueue     StorageInsertQueue
	publisher       PublishHandler
	presence        PresenceHandler
	pushHandler 	PushNotificationHandler
}

// InsertItem - Insert Item queued
type InsertItem struct {
	Event *ChannelEvent
	AppID string
}

// StorageInsertQueue - Receives all insert requests and send them into the database
type StorageInsertQueue struct {
	insertChannel chan InsertItem
}

// StoreEvent - Append channel to insert queue
func (engine *Engine) StoreEvent(appID string, event *ChannelEvent) {
	engine.insertQueue.insertChannel <- InsertItem{
		AppID: appID,
		Event: event,
	}
}

// GetCacheStorage - Get cache storage instance
func (engine *Engine) GetCacheStorage() CacheStorage {
	return engine.cacheStorage
}

// GetDeviceRepository - Get persistent repository
func (engine *Engine) GetDeviceRepository() DeviceRepository {
	return engine.databaseStorage.GetDeviceRepository()
}

// GetChannelRepository - Get persistent repository
func (engine *Engine) GetChannelRepository() ChannelRepository {
	return engine.databaseStorage.GetChannelRepository()
}

// GetAppRepository - Get persistent repository
func (engine *Engine) GetAppRepository() AppRepository {
	return engine.databaseStorage.GetAppRepository()
}

// GetClientRepository - Get persistent repository
func (engine *Engine) GetClientRepository() ClientRepository {
	return engine.databaseStorage.GetClientRepository()
}

// GetPublisher - Get Publisher handler
func (engine *Engine) GetPublisher() PublishHandler {
	return engine.publisher
}

// GetServerID - Get Server ID
func (engine *Engine) GetServerID() string {
	return engine.serverID
}

// GetPresence - Get presence handler
func (engine *Engine) GetPresence() PresenceHandler {
	return engine.presence
}

// GetPushHandler - Get Push notification handler
func (engine *Engine) GetPushHandler() PushNotificationHandler {
	return engine.pushHandler
}

var engine *Engine = nil

// GetEngine - Get engine singleton
func GetEngine() *Engine {
	return engine
}

// InitEngine - Create new engine instance
func InitEngine(dbStorage DatabaseStorage, cacheStorage CacheStorage, publisher PublishHandler, presence PresenceHandler, pushHandler PushNotificationHandler) {
	engine = &Engine{
		serverID:        xid.New().String(),
		HubsHandler:     HubsHandler{},
		databaseStorage: dbStorage,
		cacheStorage:    cacheStorage,
		insertQueue:     StorageInsertQueue{},
		publisher:       publisher,
		presence:        presence,
		pushHandler:	 pushHandler,
	}

	engine.insertQueue.insertChannel = make(chan InsertItem, 300)

	var index = 0
	var databaseWorkers = 5
	for {

		if index > databaseWorkers {
			return
		}

		go startInsertingQueue(&engine.insertQueue, dbStorage.GetChannelRepository())
		go startInsertingQueue(&engine.insertQueue, dbStorage.GetChannelRepository())
		index++
	}
}

const (
	CacheLimit   = 70 // Amount of insert before batching
	CacheTimeout = 5 * time.Second
)

func startInsertingQueue(queue *StorageInsertQueue, repo ChannelRepository) {

	cache := make([]InsertItem, 0, CacheLimit)
	tick := time.NewTicker(CacheTimeout)

	for {
		select {
		case m := <-queue.insertChannel:
			cache = append(cache, m)

			if len(cache) < CacheLimit {
				break
			}

			// Reset the timeout ticker.
			// Otherwise we will get too many sends.
			tick.Stop()

			// Send the cached messages and reset the cache.
			if err := repo.AddChannelEvents(cache); err != nil {
				log.Printf("error: %v", err)
			}
			cache = cache[:0]

			// Recreate the ticker, so the timeout trigger
			// remains consistent.
			tick = time.NewTicker(CacheTimeout)
		case <-tick.C:
			if len(cache) == 0 {
				continue
			}

			if err := repo.AddChannelEvents(cache); err != nil {
				log.Printf("error: %v", err)
			}
			cache = cache[:0]
		}
	}
}
/*
 func startInsertingQueue(queue *StorageInsertQueue, repo ChannelRepository) {

 	for {
 		select {
 		case item, isActive := <-queue.insertChannel:
 			{

 				if !isActive {
 					return
 				}

 				if err := repo.AddChannelEvent(item.AppID, item.Event.ChannelID, item.Event); err != nil {
 					log.Printf("error: %v", err)
 				}

 			}
 		}
 	}
}
*/