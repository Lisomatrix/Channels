package main

import (
	"log"

	"github.com/lisomatrix/channels/channels"
	"github.com/lisomatrix/channels/channels/auth"
	"github.com/lisomatrix/channels/channels/storage/pgxsql"
)

func main() {

	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	config, err := channels.NewConfig("./config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	auth.SetSecret(config.JWTSecret)
	pgxsql.PGXSetConnectionParams(config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.DB)

	channels.Start(config.Server.Host, config.Server.Port)

}
