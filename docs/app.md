# App

An app is a way of separating clients and channels, any publish by the a user in App `123` won't intefere with the App `321`, they can even have the same channel names, but for **Channels** they are completely different. Also, every client has an AppID associated!

!> **Note:** Almost all requests require the AppID to be sent in the headers, except WebSocket connections that also allow in the query params.

# Creating a App

Since everything is separated by Apps, how do we create them?

It is super simple, we just need a simple POST