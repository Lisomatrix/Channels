# Channels

> A simple server that treats communication like channels

Check the docs here: https://lisomatrix.github.io/Channels

## What is Channels?

Channels is a Golang server and some client libraries, mainly for [IOS](https://github.com/Lisomatrix/ChannelsSDK_Swift), [Android](https://github.com/Lisomatrix/ChannelsSDK_Android) and [Javascript](https://github.com/Lisomatrix/ChannelsSDK_JS), that allow real time communication in channels.<br>
Think of front facing publish subscribe system.

___

## Why?

> Think of a Pub/Sub with persistence and presence

This project was created to serve the purpose of creating a front facing broker and at the same time have the ability to get previous/history data.
<br><br>
**For example:** We have a chat app, with **Channels** we receive messages in real time, but we can also open our account in another device and get the messages history along with new updates!

___

## What can it do?

- App Separation by AppID
- Publish and Subscribe Channels
- Real time events with WebSockets
- Publishing with HTTP
- Subscribe with SSE (Soon)
- HTTP Channel Events Sync (On channels with persistence)
- Channel Features:
    - Close and Open (Something like freeze)
    - Private or Open (Client needs to be added or can just subscribe)
    - Optional persistence (for later access and sync)
    - Optional client presence
        - Client joined/left
        - Client online/offline
- Multiple servers with Redis
- Pluggable parts: 
    - Database (Using PostgreSQL Currently)
    - Cache (Using Redis and creating Ledis)
    - Publisher (Using Redis)
    - Presence (Using Redis)

___

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

___ 

# App

An app is a way of separating clients and channels, any publish by the a user in App `123` won't intefere with the App `321`, they can even have the same channel names, but for **Channels** they are completely different. Also, every client has an AppID associated!

!> **Note:** Almost all requests require the AppID to be sent in the headers, except WebSocket connections that also allows in the url params.

!> **Note:** Operations on Apps require the Role `Super Admin` for creating, deleting and getting apps, on update app you can use the Role `Admin` to update app associated with the `Admin` client, or just use `Super Admin`.

---

## Creating a App

Since everything is separated by Apps, how do we create them?

It is super simple, we just need a simple `POST` on `/app`.

**Headers:**
```
Authorization: token
```
**Body:**
```json
{
    "AppID": "what any id you want",
    "Name": "Any name you want"
}
```

Now you either get `200 OK` or `409 Conflict` (in case it already exists), in case you get `500 Internal Server Error` contact me :) .

---

## Updating a App

Since app have little information, you can only update their name with a even simple `PUT` on `/app/{AppID}`

**Headers:**
```
Authorization: token
```
**Body:**
```json
{
    "Name": "Any name you want"
}
```

You should get a `200 OK` or `404 Not Found`.

---

## Deleting a App

Deleting an app, as you might expect, it's just a `DELETE` on `/app/{1234}`.

**Headers:**
```
Authorization: token
```

You should get a `200 OK` or `404 Not Found`.

---

## Getting all Apps

You also can get all apps by sending a  `GET` on `/app`.

**Headers:**
```
Authorization: token
```

**Result:**

```json
{
  "Apps": [
    {
      "AppID": "1234",
      "Name": "App1"
    },
    {
      "AppID": "12344",
      "Name": "App2"
    },
    {
      "AppID": "123",
      "Name": "App3"
    }
  ]
}
```

___


# Client

Clients, as name suggests are users that connect to **Channels**, the clients can be apps, servers and web apps. <br>
Each client can have any number of devices connected at the same time.

However, while **Channels** knows about clients, their authentication and management is up to you. **Channels** doesn't authenticate clients with their credentials, it's up to you to provide a `JWT Token` to the clients so they can connect to **Channels**.

> You can see more about authentication [here]()

!> All client requests require an `Admin` role in order to manage the clients in the `Admin` app, while `Super Admin` can change all apps.
___

## Creating a Client

So in order for **Channels** to know about your clients you need to tell it.<br>
**For Example:** You register a new client and you tell **Channels** about it.

So you can just `POST` on `/client`.

> Clients have a `extra` field where you can put any `string` value including json, you should keep it small.

**Headers:**
```
Authorization: token
AppID: appID // The appID the client will belong
```
**Body:**
```json
{
    "clientID": "what any id you want",
    "username": "Any name you want",
    "extra": "A string about your user, it can be JSON"
}
```

Now you either get `200 OK` or `409 Conflict` (in case it already exists).

___


## Updating a Client

Like we saw above, clients have little information, which you can make up using the extra field, we recommend to keep it as tiny as possible, but it's up to you.

Anyway, you can update the client `username` and `extra` by sending a `PUT` on `/client/{clientID}`.

**Headers:**
```
Authorization: token
AppID: appID // The appID the client belongs
```
**Body:**
```json
{
    "username": "Any name you want",
    "extra": "A string about your user, it can be JSON"
}
```

Now you either get `200 OK` or `404 Not Found`.


___


## Deleting a Client

Losing a client is very sad, and we are sorry for you if you are checking this part.

But you can clear your tears and do it by send a `DELETE` on `/client/{clientID}`

**Headers:**
```
Authorization: token
AppID: appID // The appID the client belongs
```

Now you always get `200 OK`, if it existed now it doesn't, if it didn't in the first place **Channels** did its best.


___

## Getting clients

You can see how much your user base has grown by asking **Channels** for clients.

You can request all clients of all apps or all clients of a app. For the first simply do not include the `AppID` header.

!> **Note:** Only a **Super Admin** can get all clients from all apps, while **Admin** can only get the clients of it's apps.

You just send a `GET` on `/client`

**Headers:**
```
Authorization: token
AppID: appID // The appID you want to check, our don't send it to get all
```

**Result:**
```json

{
  "clients": [
    {
      "ID": "123",
      "Username": "lisomatrix",
      "AppID": "123",
      "Extra": "test_extra"
    },
    {
      "ID": "55",
      "Username": "lisomatrix",
      "AppID": "123",
      "Extra": "{ \"Im_Json\": \"Hi Json, Im Dad\" }"
    }
  ]
}

```

___


# Channel

> Like the name mentions, **Channels** separates things by ... channels (surprising I know). 

A channel can be any string like `my_public_channel` or `my_not_so_private_channel`. You can have repeated channels as long as they are in different apps. Also, currently we don't support topics like MQTT with `+` and `#`, a channel is just a string, think of it like and **ID**.

Currently, channels have some featues like:

- **Persistence:** Store all events published*
- **Private:** Requires and `Admin` or `Super Admin` to join/remove them
- **Presence:** Check for client connections and clients joined/removed to/from channel.

Also you have an `extra` field just like with clients to store additional data for a channel (Yes, it can be JSON).

!> Channels Management can only be made by **Admin** for it's App or **Super Admin** for every one

___

## Creating a Channel

So, in order to start publishing and subscribing and sending events you first must define the channel where you will do it.

!> **Important:** This is the hardest request in **Channels**, make sure you understand it :) 

You just send a `POST` to `/channel`.

**Headers:**
```
Authorization: token
AppID: appID // The appID the channel will belong
```
**Body:**
```json
{
    "channelID": "what any id you want",
    "name": "Any name you want",
    "persistent": true,
    "private": false,
    "presence": true,
    "extra": "A string about the channel, it can be JSON"
}
```

And you should get a `201 Created` or `409 Conflict`.

___

## Deleting a Channel

You can delete channels, and if they are persistent every single event published to this channel will be deleted!<br> **so ... be  ... careful .. like .. VERY!**

You just send a `DELETE` to `/channel/{channelID}`.

**Headers:**
```
Authorization: token
AppID: appID // The appID the channel belongs
```

And you should get a `200 OK` or `404 Not Found`.

___

## Closing a Channel

Instead of deleting a channel, if you just want events to stop being published then closing a channel is the right choice!

You just send a `POST` to `/channel/{channelID}/close`.

**Headers:**
```
Authorization: token
AppID: appID // The appID the channel belongs
```

And you should get a `200 OK` or `404 Not Found`.
And you can't publish to the channel anymore!.

___

## Opening a Channel

If you think that a channels feels alone you can open it again any time.

You just send a `POST` to `/channel/{channelID}/open`.

**Headers:**
```
Authorization: token
AppID: appID // The appID the channel belongs
```

And you should get a `200 OK` or `404 Not Found`.
 And you can start publishing to the channel again!

___

## Getting Channels (No, no the project)

You can request all channels of all apps or all channels of a app. For the first simply do not include the `AppID` header.

!> **Note:** Only a **Super Admin** can get all channels from all apps, while **Admin** can only get the channels of it's apps.

You just send a `GET` on `/channel/open` for open channels or `/channel/private` for private channels.

!>**Note:** You can't request both public and private in one request (atleast yet).

**Headers:**
```
Authorization: token
AppID: appID // The appID you want to check, our don't send it to get all
```

**Result:**
```json

{
  "channels": [
    {
      "id": "",
      "appID": "123",
      "name": "Hi_Im_Channel_One",
      "createdAt": 1615244965,
      "isClosed": false,
      "extra": "no extra",
      "isPersistent": true,
      "isPrivate": false,
      "isPresence": false
    },
    {
      "id": "1234",
      "appID": "123",
      "name": "Hi_Channel_One_Im_Dad",
      "createdAt": 1615734846,
      "isClosed": false,
      "extra": "no extra",
      "isPersistent": true,
      "isPrivate": true,
      "isPresence": false
    },
    {
      "id": "123",
      "appID": "123",
      "name": "Hi_Hi_Channel_One_Im_Dad_Im_Dad__OH_WAIT",
      "createdAt": 1615734842,
      "isClosed": false,
      "extra": "",
      "isPersistent": true,
      "isPrivate": false,
      "isPresence": true
    }
  ]
}

```

___

## Adding and Removing clients to/from Channel

For private channels, the `Admin` or `Super Admin` must add the users, of course this process could be made by a server with a `JWT Token` with one of those roles.

> Connected clients that are affected by this change are notified by receiving a `ADD_CHANNEL` or `REMOVE_CHANNEL` event, if a client is connected to a channel you just removed **Channels** will take care of unsubscribing the client, or if the client is connected and received access to a new channel it can subscribe right away!

> Also, on channels with presence enabled, connected clients will be notified that user has joined/left a channel.

So after that brief explanation, we can join or remove a client by send a `POST` to `/channel/{channelID}/join/{clientID}` to join a client or `/channel/{channelID}/leave/{clientID}` to remove.

**Headers:**
```
Authorization: token
AppID: appID // The appID the channel belongs
```

The response should be `200 OK` or `404 Not Found`.

___

## Publishing to channels

Publishing should mostly be done with WebSockets, but  some times we need to make a simple `POST` and a WebSocket would be overkill.

So we in order to publish and event with `HTTP` just send a `POST` to `/channel/{channelID}/publish`.

**Headers:**
```
Authorization: token
AppID: appID // The appID the channel belongs
```

**BODY**

```json

{
	"payload": "Your_Json_Payload",
	"eventType": "Any event type you want" // This can be used to know what to expect from the payload
}

```

And you get `200 OK` or `404 Not Found` or in case the channel is closed `400 Bad Request`.


___


# Synchronization

A lot of time our networks go down or our mobile device runs out of battery and we stop receiving messages from our friends including our friend **Channels** :(

But as soon we get up online we can start getting all data we lost while we were away :)<br> 

!> This only works with channels with persistence enabled

**We can do it some ways:**

- **Get last event channels:** last 200 messages for example
- **Get all events since specific time**: based on a timestamp
- **Get last events since specific time**: based on a timestamp and a given amount to retrieve
- **Get all events between to specific times**: also base on timestamps and they are inclusive

!> Before retrieving syncing data make sure you connection is established and you are subscribed so we can prevent missing events published between the time your fetched and subscribed!

!> Also, events are stored in database in batches, while they have timers there is a very small chance that you might lose some events that were delayed and only stored just right after you fetched them, **were are not sure this happens yet**, but could happen if the server or database are under higher load and they delay the storing of the messages.<br>If it is critical to receive those events you could fetch the events between the last timestamp of the fetched event and the first timestamp you get from the WebSocket.
<br><br>The events should still be stored, but not at the time you fetch them.
<br><br>
**But once again, we are not sure yet!**

___

## Get Last Events

In order get last events from a channel we just send a `GET` to `/last/{channelID}/{amount}`, the amount specified how much events you wan't to retrieve.

> This can be useful in a chat application in Browsers, you just need the last messages right away.

> Also the request it's pretty fast, it will try to fetch the events from cache first, but even when hitting the database with channel with more than **2 Million Rows** it takes up to **20 ms** returning **300** events.<br><br>
This test was made in a consumer **HDD** with **Database** and **Channels** running locally and using the [PGX](https://github.com/Lisomatrix/Channels/tree/main/channelserver/storage/pgxsql) implementation!

**Headers:**
```
Authorization: token
AppID: appID // The appID the channel belongs
```

**Result:**
```json

{
  "events": [
    {
      "senderID": "123",
      "eventType": "testing publish type",
      "payload": "can_be_json_or_not",
      "channelID": "123",
      "timestamp": 1615737318
    },
    {
      "senderID": "123",
      "eventType": "testing publish type",
      "payload": "can_be_json_or_not",
      "channelID": "123",
      "timestamp": 1615735212
    }
  ]
}

```

___

## Get Last Events Since

In order get last events since a given timestamp from a channel we just send a `GET` to `/last/{channelID}/{amount}/last/{timeStamp}`, the amount specified how much events you wan't to retrieve.

> This request is also pretty fast, actually even faster than the one before, im the same conditions it takes **~4ms** to return **500** events!

**Headers:**
```
Authorization: token
AppID: appID // The appID the channel belongs
```

**Result:**
```json

{
  "events": [
    {
      "senderID": "123",
      "eventType": "testing publish type",
      "payload": "can_be_json_or_not",
      "channelID": "123",
      "timestamp": 1615737318
    },
    {
      "senderID": "123",
      "eventType": "testing publish type",
      "payload": "can_be_json_or_not",
      "channelID": "123",
      "timestamp": 1615735212
    }
  ]
}

```

___ 

## Get Events Since

In order get  events since a given timestamp from a channel we just send a `GET` to `/c/{channelID}/sync/{timeStamp}`.

> This request performance depends on the amount of events are stored after (it is inclusive) the given timestamp, so be a bit careful with the timestamps or you could literally download an entire channel history!
<br><br>
**If that's what you want... well ... go ahead!**

**Headers:**
```
Authorization: token
AppID: appID // The appID the channel belongs
```

**Result:**
```json

{
  "events": [
    {
      "senderID": "123",
      "eventType": "testing publish type",
      "payload": "can_be_json_or_not",
      "channelID": "123",
      "timestamp": 1615737318
    },
    {
      "senderID": "123",
      "eventType": "testing publish type",
      "payload": "can_be_json_or_not",
      "channelID": "123",
      "timestamp": 1615735212
    }
  ]
}

```

___

## Get Last Events Between

In order get events between two timestamps from a channel we just send a `GET` to `/sync/{channelID}/{sinceTimestamp}/to/{upToTimestamp}`, the results are inclusive!.

> You should expect similar performance to the previous request

**Headers:**
```
Authorization: token
AppID: appID // The appID the channel belongs
```

**Result:**
```json

{
  "events": [
    {
      "senderID": "123",
      "eventType": "testing publish type",
      "payload": "can_be_json_or_not",
      "channelID": "123",
      "timestamp": 1615737318
    },
    {
      "senderID": "123",
      "eventType": "testing publish type",
      "payload": "can_be_json_or_not",
      "channelID": "123",
      "timestamp": 1615735212
    }
  ]
}

```