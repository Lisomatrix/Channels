package main

import (
	"log"
	_ "net/http/pprof"

	"github.com/Channels/Channels"
	"github.com/Channels/Channels/auth"
	"github.com/Channels/Channels/storage/pgxsql"
)

func main() {

	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	config, err := Channels.NewConfig("./config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	auth.SetSecret(config.JWTSecret)
	pgxsql.PGXSetConnectionParams(config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.DB)

	Channels.Start(config.Server.Host, config.Server.Port)

}
