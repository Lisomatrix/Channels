package main

import (
	"github.com/lisomatrix/channels/channels"
	"github.com/lisomatrix/channels/channels/auth"
	"github.com/lisomatrix/channels/channels/cache"
	"github.com/lisomatrix/channels/channels/core"
	"github.com/lisomatrix/channels/channels/presence"
	"github.com/lisomatrix/channels/channels/publisher"
	"github.com/lisomatrix/channels/channels/storage/mysql"
	"log"
)

func main() {
	// Make sure the config.yaml is at the same path as the binary
	config, err := channels.NewConfig("./config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	auth.SetSecret(config.JWTSecret)
	//pgxsql.PGXSetConnectionParams(config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.DB)


	// Initialize the engine with your components
	// If you make more implementations, please contribute!
	dbStorage := mysql.NewSQLStorageDatabase()
	core.InitEngine(dbStorage, cache.NewRedisCacheStorage(), publisher.NewRedisPublisher(), presence.NewRedisPresence())

	// Start will only bind the routes with handlers
	channels.Start("0.0.0.0", "8090")
}
