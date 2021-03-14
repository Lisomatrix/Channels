# App

An app is a way of separating clients and channels, any publish by the a user in App `123` won't intefere with the App `321`, they can even have the same channel names, but for **Channels** they are completely different. Also, every client has an AppID associated!

!> **Note:** Almost all requests require the AppID to be sent in the headers, except WebSocket connections that also allows in the url params.

!> **Note:** Operations on Apps require the Role `Super Admin` for creating, deleting and getting apps, on update app you can use the Role `Admin` to update app associated with the `Admin` client, or just use `Super Admin`.

---

# Creating a App

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

# Updating a App

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

# Deleting a App

Deleting an app, as you might expect, it's just a `DELETE` on `/app/{1234}`.

**Headers:**
```
Authorization: token
```

You should get a `200 OK` or `404 Not Found`.

---

# Getting all Apps

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