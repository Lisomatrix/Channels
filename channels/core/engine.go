package core

import (
	"time"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

// Engine - Holds application components
type Engine struct {
	serverID        string
	hubsHandler     *HubsHandler
	databaseStorage DatabaseStorage
	cacheStorage    CacheStorage
	insertQueue     StorageInsert
	publisher       PublishHandler
	presence        PresenceHandler
	pushHandler     PushNotificationHandler
	storageInsert   StorageInsert
	authHook        AuthHook
}

// StoreEvent - Append channel to insert queue
func (engine *Engine) StoreEvent(appID string, event *ChannelEvent) {
	if engine.storageInsert != nil {
		engine.storageInsert.StoreEvent(appID, event)
	}
}

// GetCacheStorage - Get cache storage instance
func (engine *Engine) GetCacheStorage() CacheStorage {
	return engine.cacheStorage
}
func (engine *Engine) GetHubsHandler() *HubsHandler {
	return engine.hubsHandler
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

func (engine *Engine) GetAuthHook() AuthHook {
	return engine.authHook
}

var engine *Engine = nil

// GetEngine - Get engine singleton
func GetEngine() *Engine {
	return engine
}

// EngineConfig - Config for the engine, including storage, cache, push notifications
type EngineConfig struct {
	ServerID                string                  // ServerID for server indetification, if not provided one will be generated
	HubsHandler             *HubsHandler            // If nil, then a default one is created
	DBStorage               DatabaseStorage         // Struct that holds repository
	CacheStorage            CacheStorage            // Channels, App, Sessions and events cache
	PublishHandler          PublishHandler          // Publish between servers handler
	PresenceHandler         PresenceHandler         // Handler for tracking user presence
	PushNotificationHandler PushNotificationHandler // Handler for sending push notifications
	DBWorkers               int                     // If set to -1 it will to to the default of 10
	InsertCacheLimit        int                     // Amount of events stored before batching into the database
	StorageInsert           StorageInsert           // Handler for events being stored, you can use this to batch to events, or simply ignore them. For a batching default one use StorageInsertQueue, that uses the property InsertCacheLimit
	AuthHook                AuthHook                // For the default connection, to authorize connections
}

func InitEngine(config EngineConfig) {

	SetUpLogger()

	if config.ServerID == "" {
		config.ServerID = xid.New().String()
	}

	if config.DBStorage == nil {
		log.WithFields(log.Fields{
			"EngineConfig": config,
		}).Fatal("Missing DatabaseStorage on EngineConfig")
	}

	if config.CacheStorage == nil {
		log.WithFields(log.Fields{
			"EngineConfig": config,
		}).Fatal("Missing CacheStorage on EngineConfig")
	}

	if config.PublishHandler == nil {
		log.WithFields(log.Fields{
			"EngineConfig": config,
		}).Fatal("Missing PublishHandler on EngineConfig")
	}

	if config.PresenceHandler == nil {
		log.WithFields(log.Fields{
			"EngineConfig": config,
		}).Fatal("Missing PresenceHandler on EngineConfig")
	}

	if config.PushNotificationHandler == nil {
		log.WithFields(log.Fields{
			"EngineConfig": config,
		}).Fatal("Missing PushNotificationHandler on EngineConfig")
	}

	if config.HubsHandler == nil {
		config.HubsHandler = NewHubsHandler(nil)
	}

	if config.DBWorkers == -1 {
		config.DBWorkers = 10
	}

	engine = &Engine{
		serverID:        config.ServerID,
		hubsHandler:     config.HubsHandler,
		databaseStorage: config.DBStorage,
		cacheStorage:    config.CacheStorage,
		insertQueue:     config.StorageInsert,
		publisher:       config.PublishHandler,
		presence:        config.PresenceHandler,
		pushHandler:     config.PushNotificationHandler,
		storageInsert:   config.StorageInsert,
		authHook:        config.AuthHook,
	}

	CacheLimit = config.InsertCacheLimit

	var index = 0
	for {

		if index > config.DBWorkers {
			break
		}
		go engine.storageInsert.Start(config.DBStorage.GetChannelRepository())
		index++
	}
}

var CacheLimit = 70 // Amount of insert before batching

const (
	CacheTimeout = 5 * time.Second
)
