package core

import (
	"encoding/json"
	"fmt"
	"github.com/lisomatrix/channels/channels/auth"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type getChannelMessagesRequest struct {
	ChannelID string `json:"channelID"`
	Timestamp int64  `json:"timestamp"`
}

type getChannelMessagesResponse struct {
	Events []*outEvent `json:"events"`
}

type channelPublishRequest struct {
	//ChannelID string
	Payload   string `json:"payload"`
	EventType string `json:"eventType"`
}

// CreateChannelRequest - Create channel with given ID and settings
type CreateChannelRequest struct {
	ChannelID  string   `json:"channelID"`
	Name       string   `json:"name"`
	Persistent bool     `json:"persistent"`
	Private    bool     `json:"private"`
	Presence   bool     `json:"presence"`
	Users      []string `json:"users"`
	Extra      string   `json:"extra"`
}

type outEvent struct {
	Timestamp int64  `json:"timestamp"`
	Data      []byte `json:"data"`
}

// PostEventHandler - Publish event into channel
// /channel/:channelID/publish
func PostEventHandler(context *gin.Context) {

	request := context.Request
	writer := context.Writer

	// Check for required headers
	token, appID, isOK := auth.GetAuthData(request)

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get URL channel param
	channelID, isOK := context.Params.Get("channelID")

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get request body
	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse body
	var channelPublishRequest channelPublishRequest

	err = json.Unmarshal(body, &channelPublishRequest)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check existence from cache
	channel := GetEngine().GetCacheStorage().GetChannel(appID, channelID)

	// If it doesn't exist
	// Check it from the database
	if channel == nil /*!exists*/ {
		channel, err = GetEngine().GetChannelRepository().GetAppChannel(appID, channelID)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "HTTP Publish Channel; failed to get app channel %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if channel == nil /*!exists*/ {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		GetEngine().GetCacheStorage().StoreChannel(appID, channelID, channel)
	}

	if channel.IsClosed {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	event := &ChannelEvent{
		SenderID:  identity.ClientID,
		EventType: channelPublishRequest.EventType,
		Payload:   channelPublishRequest.Payload,
		ChannelID: channelID,
		Timestamp: time.Now().Unix(),
	}

	if channel.Persistent {
		GetEngine().GetCacheStorage().StoreChannelEvent(channelID, appID, event)
		GetEngine().StoreEvent(channel.AppID, event)
	}

	// If no hub exists then we don't have clients from this hub
	hub := GetEngine().HubsHandler.ContainsHub(appID)

	if hub == nil {
		writer.WriteHeader(http.StatusOK)
		return
	}

	// If inside the hub there are no clients listening for channel, then ignore
	hubChannel := hub.ContainsChannel(channelID)

	if hubChannel == nil {
		writer.WriteHeader(http.StatusOK)
		return
	}

	// If there client listening then send to them
	hubChannel.ExternalPublish(event)

	writer.WriteHeader(http.StatusOK)
}

// CreateChannelHandler - Create channel with given info
// POST /channel
func CreateChannelHandler(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Check for required headers
	token, appID, isOK := auth.GetAuthData(request)

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get request body
	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse body
	var createChannelRequest CreateChannelRequest

	err = json.Unmarshal(body, &createChannelRequest)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	createdAt := time.Now().Unix()

	newChannel := Channel{
		ID:         createChannelRequest.ChannelID, //Not used by cache
		AppID:      appID,                          // Not used by cache
		Name:       createChannelRequest.Name,
		CreatedAt:  createdAt,
		IsClosed:   false,
		Extra:      createChannelRequest.Extra,
		Persistent: createChannelRequest.Persistent,
		Private:    createChannelRequest.Private,
		Presence:   createChannelRequest.Presence,
	}

	if isOK, err := CreateChannel(appID, &newChannel); err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Create Channel failed %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
	} else if isOK {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusConflict)
	}
}

// PostJoinChannel - Join user to a channel
// POST /channel/:channelID/join/:clientID
func PostJoinChannel(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Check for required headers
	token, appID, isOK := auth.GetAuthData(request)

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get URL channel param
	channelID, isOK := context.Params.Get("channelID")

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get URL client param
	clientID, isOK := context.Params.Get("clientID")

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	clientChannels, err := GetEngine().GetChannelRepository().GetClientAllowedChannels(clientID)

	for _, id := range clientChannels {
		if channelID == id {
			writer.WriteHeader(http.StatusConflict)
			return
		}
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Joun channel: failed to join client to channel %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isOK, err := JoinChannel(appID, channelID, clientID); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "HTTP Join channel: failed to join client %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
	} else if isOK {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}


}

// PostLeaveChannel - Remove user from channel
// POST /channel/:channelID/leave/:clientID
func PostLeaveChannel(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Check for required headers
	token, appID, isOK := auth.GetAuthData(request)

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get URL channel param
	channelID, isOK := context.Params.Get("channelID")

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get URL client param
	clientID, isOK := context.Params.Get("clientID")

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}


	clientChannels, err := GetEngine().GetChannelRepository().GetClientAllowedChannels(clientID)

	contains := false

	for _, id := range clientChannels {
		if channelID == id {
			contains = true
			break
		}
	}

	if !contains {
		writer.WriteHeader(http.StatusConflict)
		return
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Leave channel: failed to remove client from channel %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isOK, err := LeaveChannel(appID, channelID, clientID); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "HTTP Leave channel: remove client from channel %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
	} else if isOK {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}

}

// DeleteChannelHandler - Delete channel
// DELETE /channel/:channelID
func DeleteChannelHandler(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Check for required headers
	token, appID, isOK := auth.GetAuthData(request)

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get URL channel param
	channelID, isOK := context.Params.Get("channelID")

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if isOK, err := DeleteChannel(appID, channelID); err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Delete channel: failed to delete channel %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
	} else if isOK {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

// PostCloseChannel - Close channel
// POST /channel/:channelID/close
func PostCloseChannel(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Check for required headers
	token, appID, isOK := auth.GetAuthData(request)

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get URL channel param
	channelID, isOK := context.Params.Get("channelID")

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if isOK, err := SetChannelCloseStatus(appID, channelID, true); err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Close channel: failed to close channel %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
	} else if isOK {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

// PostOpenChannel - Open channel
// POST /channel/:channelID/open
func PostOpenChannel(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Check for required headers
	token, appID, isOK := auth.GetAuthData(request)

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get URL channel param
	channelID, isOK := context.Params.Get("channelID")

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if isOK, err := SetChannelCloseStatus(appID, channelID, false); err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Open channel: failed to open channel %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
	} else if isOK {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}

}

// GetChannelsResponse - Channels response data holder
type GetChannelsResponse struct {
	Channels []*Channel `json:"channels"`
}

// GetOpenChannels - Get all public channels
// GET /channel/open
func GetOpenChannels(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Get AppID and Token
	appID := request.Header.Get("AppID")
	token := request.Header.Get("Authorization")

	if token == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	identity, isOK := auth.VerifyToken(token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	var response GetChannelsResponse

	// If AppID given, and user is authorized to use AppID
	// Then fetch App specific channels
	if appID != "" && identity.CanUseAppID(appID) {

		if channels, err := GetEngine().GetChannelRepository().GetAppPublicChannels(appID); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP Get open channels: failed to get app open channels %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
		} else {
			response.Channels = channels
		}

		// Else if no given appID and user is SuperAdmin
		// Fetch all open channels from all apps
	} else if appID == "" && identity.IsSuperAdmin() {

		if channels, err := GetEngine().GetChannelRepository().GetAllPublicChannels(); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP Get open channels: failed to get all open channels %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
		} else {
			response.Channels = channels
		}
		// Otherwise there are no permissions for this request
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := json.Marshal(response)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get open channels: failed to marshal response %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}

// GetPrivateChannels - Get all private channels
// GET /channel/private
func GetPrivateChannels(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Get AppID and Token
	appID := request.Header.Get("AppID")
	token := request.Header.Get("Authorization")

	if token == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	identity, isOK := auth.VerifyToken(token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	var response GetChannelsResponse

	// If AppID given, and user is authorized to use AppID, then check if user is admin
	// Then fetch App specific channels
	if appID != "" && identity.CanUseAppID(appID) && identity.IsAdminKind() {

		if channels, err := GetEngine().GetChannelRepository().GetAppPrivateChannels(appID); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP Get open channels: failed to get app open channels %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
		} else {
			response.Channels = channels
		}

		// If AppID is given, and user is authorized to use AppID
		// Then check if is client
	} else if appID != "" && identity.CanUseAppID(appID) && identity.IsClient() {

		if channels, err := GetEngine().GetChannelRepository().GetClientPrivateChannels(identity.ClientID); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP Get open channels: failed to get app open channels %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
		} else {
			response.Channels = channels
		}

		// Else if no given appID and user is SuperAdmin
		// Fetch all open channels from all apps
	} else if appID == "" && identity.IsSuperAdmin() {

		if channels, err := GetEngine().GetChannelRepository().GetAllPrivateChannels(); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP Get open channels: failed to get all open channels %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
		} else {
			response.Channels = channels
		}
		// Otherwise there are no permissions for this request
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := json.Marshal(response)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get open channels: failed to marshal response %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}
