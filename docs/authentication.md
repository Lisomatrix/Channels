# Authentication

First, the authentication is done by your application, **Channels** only stores client usernames and IDs (more information on Creating Client Section).<br><br>After authenticating a user, you must provide him a `JWT Token`, we'll see how they should look like in a second.

## Roles

The **Channels** has three roles being these:

- `Super Admin`: can do anything since he is super!
- `Admin`: can do anything in it's own app.
- `Client`: can publish and subscribe to allowed channels on their associated App

The attribution of the roles must be made by your server, app or manually. The roles are provided by the `JWT Token`, so a client can be made a `Admin` by having a token saying so!<br>

## Token Data

For `Admin` and `Client` we need their `Role`, `ClientID` and `AppID`, for `Super Admin` we need `Role` and `ClientID`.

!> **Important:** Although we say `Super Admin` without the `_` the value in token must be `Super_Admin`


```json

{
    "Role": "Admin", // Here you can change to Client or Super_Admin
    "ClientID": "123",
    "AppID": "123" // If the role is Super_Admin we don't need this field
}

```