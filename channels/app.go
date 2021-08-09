// It provides some functions to start the channels parts
// The first option is you configure and start the engine and we bind the http routes
// The second we you configure the default parts and we start the engine and bind routes
package channels

import (
	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/gzip"
	"github.com/lisomatrix/channels/channels/connection"
	"github.com/lisomatrix/channels/channels/core"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
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

// Start channel server, make sure you configured the Engine first
func Start(host string, port string, router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)

	router.Use(CORSMiddleware())

	// WebSocket route
	//router.GET("/", connection.RequestHandler)
	router.GET("/optimized", connection.OptimizedRequestHandler)
	// router.GET("/optimized", wsHandler)

	// Only enabled GZIP Compressesion on non websocket connections
	router.Use(gzip.Gzip(gzip.DefaultCompression))

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
	//router.GET("/channel/open", core.GetOpenChannels)
	//router.GET("/channel/private", core.GetPrivateChannels)

	// Channel Sync routes
	router.GET("/sync/:channelID/:firstTimeStamp/to/:secondTimeStamp", core.GetMessagesBetweenTimeStamps)
	router.GET("/c/:channelID/sync/:lastTimeStamp", core.GetMessagesSinceTimeStamp)
	router.GET("/last/:channelID/:amount", core.GetLastMessages)
	router.GET("/last/:channelID/:amount/last/:lastTimeStamp", core.GetLastMessagesSinceTimeStamp)
	router.GET("/last/:channelID/:amount/before/:lastTimeStamp", core.GetLastMessagesBeforeTimeStamp)

	// Channel Publish
	router.POST("/channel/:channelID/publish", core.PostEventHandler)

	// Client routes
	router.POST("/client", core.CreateClientHandler)
	router.DELETE("/client/:clientID", core.DeleteClientHandler)
	router.PUT("/client/:clientID", core.UpdateClientHandler)
	router.GET("/client", core.GetClients)
	router.GET("/client/:clientID", core.GetClientHandler)

	// Presence routes
	//router.GET("/presence/:clientID", handlers.GetClientDevicesPresences)
	//router.GET("/online/:clientID", handlers.GetClientOnlineDevices)

	log.Info("Running on host %s and port %v", host, port)
	log.Fatal(router.Run(host + ":" + port))
}
