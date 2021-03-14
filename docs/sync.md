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