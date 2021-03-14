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