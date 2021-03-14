package main

import (
	"log"
	_ "net/http/pprof"

	"lisomatrix.pt/channelserver/channelserver"
	"lisomatrix.pt/channelserver/channelserver/auth"
	"lisomatrix.pt/channelserver/channelserver/storage/pgxsql"
)

func main() {

	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	config, err := NewConfig("./config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	auth.SetSecret(config.JWTSecret)
	pgxsql.PGXSetConnectionParams(config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.DB)

	channelserver.Start(config.Server.Host, config.Server.Port)

}
