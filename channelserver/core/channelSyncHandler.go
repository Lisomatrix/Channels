package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"lisomatrix.pt/channelserver/channelserver/auth"
)

type getChannelEventsResponse struct {
	Events []*ChannelEvent `json:"events"`
}

// GetMessagesBetweenTimeStamps - Fetch messages between timestamps
// GET /channel/:channelID/sync/:firstTimeStamp/to/:secondTimeStamp
func GetMessagesBetweenTimeStamps(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Check for required headers
	token, appID, isOK := auth.GetAuthData(request)

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate token
	identity, isOK := auth.VerifyToken(token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	channelID := context.Params.ByName("channelID")

	if channelID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	firstTimeStampStr := context.Params.ByName("firstTimeStamp")
	secondTimeStampStr := context.Params.ByName("secondTimeStamp")

	if firstTimeStampStr == "" || secondTimeStampStr == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	firstTimeStamp, err := strconv.ParseInt(firstTimeStampStr, 10, 64)
	secondTimeStamp, err := strconv.ParseInt(secondTimeStampStr, 10, 64)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get messages between timestamps: failed convert timestamp %v\n", err)
	}

	// Check if channel exists
	exists, err := GetEngine().GetChannelRepository().ExistsAppChannel(appID, channelID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get messages between timestamps: failed to check app channel existence %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// Get events
	events, err := GetEngine().GetChannelRepository().GetChannelEventsAfterAndBefore(appID, channelID, firstTimeStamp, secondTimeStamp)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get messages between timestamps: failed fetch events %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := getChannelEventsResponse{Events: events}

	data, err := json.Marshal(response)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get messages since timestamp: failed to marshal response %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}

// GetMessagesSinceTimeStamp - Fetch messages after timestamp
// GET /channel/:channelID/sync/:lastTimeStamp
func GetMessagesSinceTimeStamp(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Check for required headers
	token, appID, isOK := auth.GetAuthData(request)

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate token
	identity, isOK := auth.VerifyToken(token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	channelID := context.Params.ByName("channelID")

	if channelID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	lastTimeStampStr := context.Params.ByName("lastTimeStamp")

	if lastTimeStampStr == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	lastTimeStamp, err := strconv.ParseInt(lastTimeStampStr, 10, 64)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get messages since timestamp: failed convert timestamp %v\n", err)
	}

	// Check if channel exists
	exists, err := GetEngine().GetChannelRepository().ExistsAppChannel(appID, channelID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get messages since timestamp: failed to check app channel existence %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// Get events
	events, err := GetEngine().GetChannelRepository().GetChannelEventsAfter(appID, channelID, lastTimeStamp)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get messages since timestamp: failed fetch events %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := getChannelEventsResponse{Events: events}

	data, err := json.Marshal(response)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get messages since timestamp: failed to marshal response %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}

// GetLastMessagesSinceTimeStamp - Fetch last messages after a timestamp
// GET /channel/:channelID/sync/:lastTimeStamp/last/:amount
func GetLastMessagesSinceTimeStamp(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Check for required headers
	token, appID, isOK := auth.GetAuthData(request)

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate token
	identity, isOK := auth.VerifyToken(token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	channelID := context.Params.ByName("channelID")

	if channelID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	amountStr := context.Params.ByName("amount")

	if amountStr == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get last messages since timestamp: failed convert amount %v\n", err)
	}

	lastTimeStampStr := context.Params.ByName("lastTimeStamp")

	if lastTimeStampStr == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	lastTimeStamp, err := strconv.ParseInt(lastTimeStampStr, 10, 64)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get last messages since timestamp: failed convert timestamp %v\n", err)
	}

	// Check if channel exists
	exists, err := GetEngine().GetChannelRepository().ExistsAppChannel(appID, channelID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get last messages since timestamp: failed to check app channel existence %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	events, err := GetEngine().GetChannelRepository().GetChannelLastEventsAfter(appID, channelID, amount, lastTimeStamp)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get last messages since timestamp: failed fetch events %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := getChannelEventsResponse{Events: events}

	data, err := json.Marshal(response)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get last messages since timestamp: failed to marshal response %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}

// GetLastMessages - Fetch last messages
// GET /channel/:channelID/last/:amount
func GetLastMessages(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Check for required headers
	token, appID, isOK := auth.GetAuthData(request)

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate token
	identity, isOK := auth.VerifyToken(token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	channelID := context.Params.ByName("channelID")

	if channelID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	amountStr := context.Params.ByName("amount")

	if amountStr == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get last messages: failed convert amount %v\n", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if channel exists
	exists, err := GetEngine().GetChannelRepository().ExistsAppChannel(appID, channelID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get last messages: failed to check app channel existence %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	var events []*ChannelEvent

	if amount <= CacheQueueSize {
		size := GetEngine().GetCacheStorage().GetChannelEventsSize(channelID)

		if size >= uint64(amount) {
			events = GetEngine().GetCacheStorage().GetChannelEvents(channelID, amount)
		} else {
			events, err = GetEngine().GetChannelRepository().GetChannelLastEvents(appID, channelID, amount)

			if err != nil {
				fmt.Fprintf(os.Stderr, "HTTP Get last messages: failed fetch events %v\n", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

	} else {
		// Get events
		events, err = GetEngine().GetChannelRepository().GetChannelLastEvents(appID, channelID, amount)

		if err != nil {
			fmt.Fprintf(os.Stderr, "HTTP Get last messages: failed fetch events %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Prepare response
	response := getChannelEventsResponse{Events: events}

	data, err := json.Marshal(response)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get last messages: failed to marshal response %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}
