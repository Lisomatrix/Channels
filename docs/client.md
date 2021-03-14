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