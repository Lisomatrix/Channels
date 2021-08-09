package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lisomatrix/channels/channels"
	"github.com/lisomatrix/channels/channels/auth"
	"github.com/lisomatrix/channels/channels/cache"
	"github.com/lisomatrix/channels/channels/core"
	"github.com/lisomatrix/channels/channels/presence"
	"github.com/lisomatrix/channels/channels/publisher"
	"github.com/lisomatrix/channels/channels/push"
	"github.com/lisomatrix/channels/channels/storage/gormsql"
	"github.com/lisomatrix/channels/channels/storage/pgxsql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PublishHandler struct{}

func (h *PublishHandler) OnPublish(channelID string, channelEvent *core.ChannelEvent, shouldStore bool) bool {

	return false
}

func (h *PublishHandler) CanPublish(channelID string, channelEvent *core.ChannelEvent, shouldStore bool, session *core.Session, publishRequest *core.PublishRequest) bool {

	return false
}

//func (h *PublishHandler)

type SimpleAuthHook struct{}

func (h *SimpleAuthHook) Authenticate(token, appID, deviceID string, request *http.Request) *auth.Identity {
	identity := auth.Identity{
		Role:     "Client",
		AppID:    appID,
		ClientID: "lisomatrix",
	}

	return &identity
}

type SimpleSessionHook struct{}

func (h *SimpleSessionHook) OnInitialized(session *core.Session) {
	log.Println("Session initialized")
}

func (h *SimpleSessionHook) OnClose(session *core.Session) {
	log.Println("Session closed")
}

func (h *SimpleSessionHook) CanSubscribe(channelID string, session *core.Session, isAllowedChannel bool) bool {
	return true
}

func (h *SimpleSessionHook) CanPublish(channelID string, session *core.Session, isAllowedChannel bool) bool {
	return true
}

type HubHook struct{}

func (h *HubHook) OnClose(hub *core.Hub)                            {}
func (h *HubHook) OnChannelRemoved(channelID string, hub *core.Hub) {}
func (h *HubHook) OnSessionAdded(session *core.Session, hub *core.Hub) {
	session.SetHook(&SimpleSessionHook{})
}
func (h *HubHook) OnSessionRemoved(session *core.Session, hub *core.Hub) {}
func (h *HubHook) OnPublish(channelID string, channelEvent *core.ChannelEvent, shouldStore bool, session *core.Session) (bool, bool) {

	channelEvent.Payload = "Message was intercepted"

	if err := session.Send(channelEvent); err != nil {
		log.Println(err)
	}

	if !shouldStore {
		shouldStore = true
	}

	return true, shouldStore
}
func (h *HubHook) OnSubscribe(channelID string, session *core.Session) bool {
	return true
}
func (h *HubHook) OnUnsubscribe(channelID string, session *core.Session) {}

type HubsHandlerHook struct{}

func (h *HubsHandlerHook) OnNewHub(hub *core.Hub) core.HubHook {
	log.Println("Hub initialized")
	return &HubHook{}
}

func (h *HubsHandlerHook) OnRemoveHub(hub *core.Hub) {
	log.Println("Hub removed")
}

func main() {
	// configYa, err := channels.NewConfig("./config.yaml")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	auth.SetSecret("tictic5532")
	//configYa.Database.User
	pgxsql.PGXSetConnectionParams("postgres", "postgres", "127.0.0.1", "5432", "postgres")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	stor := gormsql.NewGormDatabaseStorage(db)

	//dbStorage := pgxsql.NewSQLStorageDatabase()
	cacheHandler := cache.NewLedisCacheStorage()
	presenceHandler := presence.NewLedisPresenceWithDB(cacheHandler.GetDB())
	publisherHandler := publisher.NewEmptyPublisher()
	pushHandler := &push.EmptyPushNotificationHandler{}

	config := core.EngineConfig{
		ServerID:    "",
		HubsHandler: core.NewHubsHandler(&HubsHandlerHook{}),
		DBStorage:   stor,
		//DBStorage:               dbStorage,
		CacheStorage:            cacheHandler,
		PublishHandler:          publisherHandler,
		PresenceHandler:         presenceHandler,
		PushNotificationHandler: pushHandler,
		DBWorkers:               -1,
		InsertCacheLimit:        70,
		StorageInsert:           core.NewStorageInsertQueue(),
		AuthHook:                &SimpleAuthHook{},
	}

	core.InitEngine(config)

	gin.DisableConsoleColor()

	// core.CreateChannel("lisomatrix", &core.Channel{
	// 	ID:         xid.New().String(),
	// 	AppID:      "lisomatrix",
	// 	Name:       "Lisomatrix Channel",
	// 	CreatedAt:  time.Now().Unix(),
	// 	IsClosed:   false,
	// 	Extra:      "",
	// 	Persistent: false,
	// 	Private:    false,
	// 	Presence:   false,
	// 	Push:       false,
	// })

	router := gin.Default()

	channels.Start("0.0.0.0", "80", router)
}
