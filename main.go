package main

import (
	"log"
	_ "net/http/pprof"

	"github.com/channelserver/channelserver"
	"github.com/channelserver/channelserver/auth"
	"github.com/channelserver/channelserver/storage/pgxsql"
)

func main() {

	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	config, err := channelserver.NewConfig("./config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	auth.SetSecret(config.JWTSecret)
	pgxsql.PGXSetConnectionParams(config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.DB)

	channelserver.Start(config.Server.Host, config.Server.Port)

}
