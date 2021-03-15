# JavaScript SDK

If you want to connect to **channels** with your browser, then this is the right place for you!<br>
First get the code from [here](https://github.com/Lisomatrix/ChannelsSDK_JS) and check the authentication section for creating a `JWT Token`.<br>

After getting the token, initialize the SDK with:

```javascript
let channelsSDK = new ChannelsSDK({ 
    url: '://url:port', // Don't fill behind the ://
    appID: 'AppID',
    token: 'JWT Token',
    secure: false // If should be WSS or WS and HTTPS or HTTP
});
```
> You need to keep the object around and using for almost everything
___

## Getting your channels

To get the client channels you can retrieve the channels that are public or the ones that are private. To do that just use the following:

> The package works in browser, it uses the native WebSocket, if you can replace it, then it can work on Node

```javascript
channelsSDK.fetchPublicChannels().then(channels => { // Replace public with private for private channels
});
```

## Listening for new channels and removed channels

You can know you when lost access to a channel for received with:

```javascript
channelsSDK.setOnChannelAdded((channelID) => {
    // When you receive access to a channel
})

channelsSDK.setOnChannelRemoved((channelID) => {
    // When you lose access to a channel
})
```

___

## Working with a Channel

The object `Channel` is the object you will use the most, with it you can subscribe, publish, get other clients presence and get events.<br>

First, in order to get a instance you can get from `getPublicChannels()` or `getPrivateChannels()` or after you get the public or private channels you can get one with `channelsSDK.getChannel("channelID");`, **if none is found it will return null!**

Once you have a `Channel` you subscribe with:

```javascript
channel.subscribe(() => {
    // Callback for when the subscribe is confirmed
})
```

And you can get events with:

```javascript
channel.setOnMessage((event) => {
    // Callback for the event
});
```

For presences events you have:

```
// On initial status from the channel
channel.setOnInitialStatusUpdate(() => {
    // You can get the presence with
    channel.getPresencesStatus();
});

// When status changes
channel.setOnOnlineStatusUpdate((statusUpdate) => {
    // The user status that changed
});

// When a user is added to the channel
channel.setOnJoin((join) => {

});

// When a user is removed from the channel
channel.setOnLeave((leave) => {

});
```

___

## Publishing

To publish messages is simple as:

```javascript

    channel.publish("my_event", "my_payload_can_be_json", () => {
        // Published confirmation
    });

    // If you don't want the confirmation or don't want the event to be stored pass null on the callback
    channel.publish("my_event", "my_payload_can_be_json", null);

```

!> **Important:** If your event doesn't request a confirmation, `Channels` will consider that the event is not important to store on channels with persistence enabled!


___

## Getting last events

For getting events it's as simple as:

> You can check other types of getting events on the section **Synchronization**

```javascript

// Get last X (5 in this case) events 
channel.fetchLastEvents(5).then(events => {
    console.log(events);
});

// Get X (5 in this case) since given timestamp (inclusive)
channel.fetchLastEventsSince(5, 1615570823).then(events => {
    console.log(events);
});

// Get events since timestamp (inclusive)
channel.fetchEventsSince(1615570823).then(events => {
    console.log(events);
});

// Since timestamp to Up to timestamp (inclusive)
channel.fetchEventsBetween(1615570823, 1615570823).then(events => {
    console.log(events);
});

```