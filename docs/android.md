# Android SDK

If you want to connect to **channels** with your devices, then this is the right place for you!<br>
First get the code from [here](https://github.com/Lisomatrix/ChannelsSDK_Android) and check the authentication section for creating a `JWT Token`.<br>

After getting the token, initialize the SDK with:

```java
ChannelsSDK.initialize(context, "wss://server_url:port","AppID", "your_jwt_Token");
```

This should connect to the server if both ulr and credentials are right.


___

## Getting your channels

To get the client channels you can retrieve the channels that are public or the ones that are private. To do that just use the following:

```java
ChannelsSDK
        .getInstance()
        .getChannelService()
        .getPublicChannels(new GetChannelsCallback() { // For Private just replace the word public
            @Override
            public void onSuccess(List<ChannelInfo> channelInfos) {
                // Here you get your channels information
            }

            @Override
            public void onError(Throwable throwable) {
                // In case an error happens
            }
        });
```

With the `ChannelInfo` you can get a `Channel`, with the static method `Channel.fromChannelInfo(ChannelInfo)`;<br>
We will see what you can do with it in a second.

___

## Listening for new channels and removed channels

You can know you when lost access to a channel for received with:

```java
ChannelsSDK.setChannelsListener(new ChannelsListener() {
            @Override
            public void onChannelRemoved(String channelID) {
                // When your remove a client from a channel
            }

            @Override
            public void onChannelAdded(String channelID) {
                // When you add a channel to a client
            }
        });
```

___

## Working with a Channel

The object `Channel` is the object you will use the most, with it you can subscribe, publish, get other clients presence and get events.<br>

First, in order to get a instance you can use the `ChannelInfo` you get from `getPublicChannels()` or `getPrivateChannels()` with `Channel.fromChannelInfo(ChannelInfo)`, or after you get the public or private channels you can get one with `ChannelsSDK.getChannel("channelID");`, **if none is found it will return null!**

!> Please don't try to create an instance out of the ways we provide, but in case you do register it with `ChannelsHandler.getInstance().registerChannel(channel);` or you won't receive any updates.

Once you have a `Channel` you subscribe with:

```java
    m_channel.subscribe(new ChannelListener() {
            @Override
            public void onPublishAcknowledge(RequestAcknowledge requestAcknowledge) {
                // When you publick with ACK enabled
            }

            @Override
            public void onSubscribed() {
                // Callback so you know when you are subscribed
            }

            @Override
            public void onChannelEvent(ChannelEvent event) {
                // When your receive a message
            }

            @Override
            public void onRemoved() {
                // When you lost access to a channel
            }
        });
```


For a `channel` with presence enabled you have another callback, a long one!

```java
m_channel.setPresenceListener(new ChannelPresenceListener() {
            @Override
            public void onClientJoinChannel(ClientJoin clientJoin) {
                // When a user is added to a channel
            }

            @Override
            public void onClientLeaveChannel(ClientLeave clientLeave) {
                // When a user is removed from a channel
            }

            @Override
            public void onOnlineStatusUpdate(OnlineStatusUpdate onlineStatusUpdate) {
                
                // You receive an event when a client goes online or offline
                // And you can get the all the presence information with:
                Map<String, ClientPresenceStatus> presenceStatusMap = m_channel.getPresences();
                // The key in the map is the ClientID
            }

            @Override
            public void onInitialStatusUpdate() {
               
               // After you subscribe to a channel, the channel should send it's
               // Current presence state to you
               // When you get it, this 'little' call is called and here you can get presences with:
               Map<String, ClientPresenceStatus> presenceStatusMap = m_channel.getPresences();
            }


});
```

> You should this set this listener before subscribing but you still can get the presence status with the `.getPresences()`.

___

## Publishing

To publish messages is simple as:

```java
m_channel.publish("event_type", "any_payload", true);
/* Set to true when you want to get a confirmation */
```

!> **Important:** If your event doesn't request a confirmation, `Channels` will consider that the event is not important to store on channels with persistence enabled!

___

## Getting last events

For getting events it's as simple as:

> You can check other types of getting events on the section **Synchronization**

```java
    ChannelsSDK.getChannelLastEvents("channelID", 20 /* Amount */, new GetChannelEventsCallback() {
            @Override
            public void onSuccess(List<ChannelEvent> events) {
                
            }

            @Override
            public void onError(Throwable throwable) {

            }
    });

    or 

    m_channel.getLastEvents(10 /* Amount */, new GetChannelEventsCallback() {
            @Override
            public void onSuccess(List<ChannelEvent> events) {
                
            }

            @Override
            public void onError(Throwable throwable) {

            }
    });
```