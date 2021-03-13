package channelserver

import (
	"fmt"
	"log"

	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"lisomatrix.pt/channelserver/channelserver/cache"
	"lisomatrix.pt/channelserver/channelserver/connection"
	"lisomatrix.pt/channelserver/channelserver/handlers"
	"lisomatrix.pt/channelserver/channelserver/presence"
	"lisomatrix.pt/channelserver/channelserver/publisher"

	"lisomatrix.pt/channelserver/channelserver/core"
	"lisomatrix.pt/channelserver/channelserver/storage/pgxsql"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, AppID")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Start channel server
func Start(host string, port string) {
	fmt.Printf("Process ID: %v \n", os.Getpid())

	// TODO: Explore message batching
	// TODO: For example, when a client send a message he waits for a reply confirmation, during this time we can chunck messages

	// TODO: Create Device HTTP handler

	// TODO: FCM implementation

	//dbStorage := storagesql.NewSQLStorageDatabase()
	dbStorage := pgxsql.NewSQLStorageDatabase()
	core.InitEngine(dbStorage, cache.NewRedisCacheStorage(), publisher.NewRedisPublisher(), presence.NewRedisPresence())

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(corsMiddleware())

	// WebSocket route
	router.GET("/", connection.RequestHandler)
	router.GET("/optimized", connection.OptimizedRequestHandler)
	// router.GET("/optimized", wsHandler)

	// Device routes
	router.POST("/device", core.CreateDevice)
	router.DELETE("/device/:deviceID", core.RemoveDevice)

	// App routes
	router.POST("/app", core.CreateApp)
	router.DELETE("/app/:appID", core.DeleteApp)
	router.PUT("/app/:appID", core.UpdateApp)
	router.GET("/app", core.GetApps)

	// Channel management routes
	router.POST("/channel", core.CreateChannelHandler)
	router.POST("/channel/:channelID/join/:clientID", core.PostJoinChannel)
	router.POST("/channel/:channelID/leave/:clientID", core.PostLeaveChannel)
	router.DELETE("/channel/:channelID", core.DeleteChannelHandler)
	router.POST("/channel/:channelID/close", core.PostCloseChannel)
	router.POST("/channel/:channelID/open", core.PostOpenChannel)
	router.GET("/channel/open", core.GetOpenChannels)
	router.GET("/channel/private", core.GetPrivateChannels)

	// Channel Sync routes
	router.GET("/sync/:channelID/:firstTimeStamp/to/:secondTimeStamp", core.GetMessagesBetweenTimeStamps)
	router.GET("/c/:channelID/sync/:lastTimeStamp", core.GetMessagesSinceTimeStamp)
	router.GET("/last/:channelID/:amount", core.GetLastMessages)
	router.GET("/last/:channelID/:amount/last/:lastTimeStamp", core.GetLastMessagesSinceTimeStamp)

	// Channel Publish
	router.POST("/channel/:channelID/publish", core.PostEventHandler)

	// Client routes
	router.POST("/client", core.CreateClientHandler)
	router.DELETE("/client/:clientID", core.DeleteClientHandler)
	router.PUT("/client/:clientID", core.UpdateClientHandler)
	router.GET("/client", core.GetClients)
	router.GET("/client/:clientID", core.GetClientHandler)

	// Presence routes
	router.GET("/presence/:clientID", handlers.GetClientDevicesPresences)
	router.GET("/online/:clientID", handlers.GetClientOnlineDevices)

	log.Println("Running on host " + host + " port: " + port)
	log.Fatal(router.Run(host + ":" + port))
}
