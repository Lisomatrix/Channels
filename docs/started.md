# Getting Started

## Easy way

The file [app.go](https://github.com/Lisomatrix/Channels/blob/main/channelserver/app.go) provides a function **Start(host string, port string)** that starts the Channel Servers with the default settings, currently the default settings are **PostgreSQL** for storage and the rest is using **Redis**, these can be changed!

Before starting we must provide the connection settings to PostgreSQL, we can do that by providing an [config.yaml](https://github.com/Lisomatrix/Channels/blob/main/example_config.yaml), and the SQL Schema [here](www.example.com).

After that make sure you have your redis running locally and the server should start!

## Bit harder way

Looking at the file [app.go](https://github.com/Lisomatrix/Channels/blob/main/channelserver/app.go), we see that we need instances of the structs that implement the following interfaces:

- [Storage interfaces](https://github.com/Lisomatrix/Channels/blob/main/channelserver/core/storage.go)

- [Presence interface](https://github.com/Lisomatrix/Channels/blob/main/channelserver/core/presence.go)

- [Publisher interface](https://github.com/Lisomatrix/Channels/blob/main/channelserver/core/publisher.go)

- [Cache interface](https://github.com/Lisomatrix/Channels/blob/main/channelserver/core/cache.go)

We currently have [PostgresSQL Storage implementation](https://github.com/Lisomatrix/Channels/blob/main/channelserver/storage/pgxsql/pgxStorage.go) (check [here](https://github.com/Lisomatrix/Channels/tree/main/channelserver/storage/storagesql) for a database/sql implementation). We also have [Redis Presence implementation](https://github.com/Lisomatrix/Channels/blob/main/channelserver/presence/redisPresence.go), [Redis Publisher implementation](https://github.com/Lisomatrix/Channels/tree/main/channelserver/publisher) and [Redis Cache implementation](https://github.com/Lisomatrix/Channels/blob/main/channelserver/cache/redisCache.go), a Ledis cache implementation is in the works!

Looking again at [app.go](https://github.com/Lisomatrix/Channels/blob/main/channelserver/app.go), we just need to initialize the **Engine**, call **core.InitEngine(storage, cache, publisher, presence)**, and now you can use **core.Engine** for the Channels main logic, the object is accessible everywhere with **core.GetEngine()** and holds the interfaces provided at init.

In case you pretend to make your own HTTP handlers or some custom logic you can use some helpers like this [Channel Helper](https://github.com/Lisomatrix/Channels/blob/main/channelserver/core/channelHelper.go), [Client Helper](https://github.com/Lisomatrix/Channels/blob/main/channelserver/core/clientHelper.go) and [Hubs Handler](https://github.com/Lisomatrix/Channels/blob/main/channelserver/core/hubsHandler.go) (this one can be accessed with **core.GetEngine().HubsHandler**) to avoid repeating yourself.