package main

import (
	"log"

	"github.com/lisomatrix/channels/channels"
	"github.com/lisomatrix/channels/channels/auth"
	"github.com/lisomatrix/channels/channels/storage/pgxsql"
)

func main() {

	// Make sure the config.yaml is at the same path as the binary
	config, err := channels.NewConfig("./config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	// Since the components are pluggable, you must initialize them before
	// In this case we are using the default ones, so just set their properties

	// The InitEngineAndStart function initializes the engine and binds routes
	auth.SetSecret(config.JWTSecret)
	pgxsql.PGXSetConnectionParams(config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.DB)

	// Initializes the engine and starts
	channels.InitEngineAndStart("0.0.0.0", "8090")
}
