# IOS SDK

If you want to connect to **channels** with your devices, then this is the right place for you!<br>
First get the code from [here](https://github.com/Lisomatrix/ChannelsSDK_Swift) and check the authentication section for creating a `JWT Token`.<br>

After getting the token, initialize the SDK with:

```swift
// Not for the URL don't put schema just for example 192.168.1.2 or example.pt
ChannelsSDK.shared().initialize(url: "example.com:port", appID: "your app ID", secure: false)
// After initializing you can connect with
ChannelsSDK.shared().connect(token: "your jwt token", deviceID: nil) // nil deviceID to auto generate
```

This should connect to the server if both ulr and credentials are right.


___

## Getting your channels

To get the client channels you can retrieve the channels that are public or the ones that are private. To do that just use the following:

```swift
let channelsAPI = ChannelsSDK.shared().getChannelsAPI()

channelsAPI.getPublicChannels() { channels, isOK in // Change the Public to private for private channels
            if (isOK) {
                print("Public \\(channels!)")
            }
        }
```
___

## Listening for new channels and removed channels

You can know you when lost access to a channel for received by implementing `ChannelsAccessListener` protocol and:

```swift

// Set listener
ChannelsSDK.shared().setChannelAccessListener(listener: self)

// And you get your events
func onChannelAdded(channelID: String) {
    // When you add a channel to a client
}
    
func onChannelRemoved(channelID: String) {
    // When your remove a client from a channel
}
```

___

## Working with a Channel

The object `Channel` is the object you will use the most, with it you can subscribe, publish, get other clients presence and get events.<br>

First, in order to get a instance you can get from `getPublicChannels()` or `getPrivateChannels()` or after you get the public or private channels you can get one with `ChannelsSDK.shared().getChannel(channelID: "channelID");`, **if none is found it will return null (I mean NIL!)!**

Once you have a `Channel` you subscribe by implementing `ChannelListener` and:

```swift
    channel.subscribe(listener: self)

    // And you get your events

    func onPublishAcknowledge(ack: RequestAcknowledge) {
        /* Confirmation that publish was received */
    }
    
    func onSubscribed() {
        /* Confirmation that you subscribed */
    }
    
    func onChannelEvent(event: EventChannel) {
        /* Handle new channel events */
        print("Event received of type \(event.eventType) with data \(event.payload)");
    }
    
    func onRemoved() {
        /* You lost access to channel */
    }
```


For a `channel` with presence enabled you have another callback, a long one!

```swift
    func onClientJoinChannel(clientJoin: ClientJoined) {
        /* When client is added to channel */
    }
    
    func onClientLeaveChannel(clientLeave: ClientLeft) {
        /* When client is removed from channel */
    }
    
    func onOnlineStatusUpdate(onlineStatusUpdate: ClientPresenceStatus) {
        /* When client goes offline or online */
    }
    
    func onInitialStatusUpdate() {
        /* When you received the initial presence state from a channel when you subscribe */
        let presences: Dictionary<String, ClientPresenceStatus> = channel!.getPresences()
        print(presences["clientID"]?.status) // IsOnline
        print(presences["clientID"]?.timestamp) // last timestamp
    }
```

> You should this set this listener before subscribing but you still can get the presence status with the `.getPresences()`.

___

## Publishing

To publish messages is simple as:

```swift
channel.publish(eventType: "event_type", payload: "you_payload_can_be_json", notify: true)
/* Set to true when you want to get a confirmation */
```

!> **Important:** If your event doesn't request a confirmation, `Channels` will consider that the event is not important to store on channels with persistence enabled!

___

## Getting last events

For getting events it's as simple as:

> You can check other types of getting events on the section **Synchronization**

```swift
    let channelsAPI = ChannelsSDK.shared().getChannelsAPI()

    channelsAPI.getLastChannelEvents(channelID: "channelID", amount: 10) { events, isOK in
            if (isOK) {
                print("Events: \\(events!)")
            }
    }

    or 

    chann.getLastEvents(amount: 50) { events, isOK in
            if (isOK) {
                print("Events: \\(events!)")
            }
    }
```