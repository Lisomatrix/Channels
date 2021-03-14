# Multiple Servers

**Channels** can be used with multiple servers using a Pub/Sub system... wait ... ain't this Pub/Sub already?<br>
Well, you are right, but pushing the work to a already battle tested and lightweight software like **Redis** works to, and very well acctually, and saves us time for other things, and we can also use it for caching!<br><br>
But, if you don't like the solution and want to implement other thing other than **Redis** go ahead! You just need to implement [this](https://github.com/Lisomatrix/Channels/blob/main/channelserver/core/publisher.go) interface and you are good to go!

## But how this works exactly?

Well, it's acctually pretty simple.<br><br>
First let's describe how **Channels** works with apps and channels.<br>
When a new user connects to **Channels** it checks if the `Hub` of the app, which the client is associated, is currently active in that server, if not it creates a `Hub` instance.

> Think of a `Hub` like a class that routes messages to channels and subscribes clients to channels. Each `Hub` represents an app, and we can have lots of apps in one server and just one in another. They are created as needed, something along the lines of lazy loading.

Inside each `Hub` we can have a lot of `channels` or maybe none. Like the the app they are loaded as needed. Once a `channel` is loaded `Channels` will subscribe to a topic like `{AppID}:{ChannelID}` in `Redis` (with the current implementation) and start to send to the subscribed client the external data they are getting.<br>
For presence information updates, pub/sub is used, but also the K/V store in `Redis` to store information and use it as shared data, this helps to detect if a user is offline by checking if there aren't any devices associated with a user that are still connected.<br><br>


Another question you might have is, what happens if **Redis** goes down? Well, some events could not be broadcasted to other instances, and presence information would be incorrect, but you could also try redis clustering, or bring your publisher implementation.<br><br>
And what happens if we lose access to database? In that case I'm sorry, but for now your events won't be stored, if we stored them in memory the server could go down pretty fast with lack of memory.

> Event if both things fail, **Channels** should not crash, but won't have the desired results either.