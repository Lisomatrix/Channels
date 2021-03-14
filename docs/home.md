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
