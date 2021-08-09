package core

import (
	"net/http"

	"github.com/lisomatrix/channels/channels/auth"
)

type HubsHandlerHook interface {
	OnNewHub(hub *Hub) HubHook
	OnRemoveHub(hub *Hub)
}

type HubHook interface {
	// Called when there are no more sessions on this hub, then it is removed to save memory
	OnClose(hub *Hub)

	// When a channels is deleted/unused for some time or closed
	// It closes all connections and then calls this function
	OnChannelRemoved(channelID string, hub *Hub)

	// Called when a new session connects, you may use it to set a SessionHook
	OnSessionAdded(session *Session, hub *Hub)
	// Called when a new session disconnects
	OnSessionRemoved(session *Session, hub *Hub)

	// Called before publishing, you may return false to prevent from publishing
	// Also you may set shouldStore to false in order to prevent the event from being stored on the DB
	// The first return param is if publish should be cancelled or not, the second is if the event should be stored or not
	// For default behaviour return the shouldStore property
	OnPublish(channelID string, channelEvent *ChannelEvent, shouldStore bool, session *Session) (bool, bool)
	// Called before subscribing, you may return false to prevent the session from subscribing
	OnSubscribe(channelID string, session *Session) bool
	// Called after a session unsubscribe
	OnUnsubscribe(channelID string, session *Session)
}

type SessionHook interface {
	OnInitialized(session *Session)
	OnClose(session *Session)

	// Called while checking if user can subscribe, the isAllowedChannel means if the channel is in the allowed list
	// You must return a bool if the user can subscribe or not, you may return isAllowedChannel for the default behaviour
	CanSubscribe(channelID string, session *Session, isAllowedChannel bool) bool

	// Called while checking if use can publish, the isAllowedChannel means if the channel is in the allowed list
	// You must return a bool if the user can publish or not, you may return isAllowedChannel for the default behaviour
	CanPublish(channelID string, session *Session, isAllowedChannel bool) bool
}

type AuthHook interface {
	Authenticate(token, appID, deviceID string, request *http.Request) *auth.Identity
}
