package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Channels/Channels/auth"
	"github.com/Channels/Channels/core"
	"github.com/gin-gonic/gin"
)

type lastDeviceStatusResponse struct {
	Devices []*core.LastDevicePresence `json:"devices"`
}

type onlineDevicesResponse struct {
	OnlineDeviceIDs []string `json:"onlineDevices"`
}

// GetClientOnlineDevices - Get all client online devices
// GET /online/:clientID
func GetClientOnlineDevices(context *gin.Context) {
	writer := context.Writer
	request := context.Request

	tokenString := request.Header.Get("Authorization")

	// Token must be present
	if tokenString == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	identity, isOK := auth.VerifyToken(tokenString)

	// Invalid token
	if !isOK {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// If not super admin, then ther client doesn't have enough permissions
	if !identity.IsSuperAdmin() {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	clientID := context.Params.ByName("clientID")

	if clientID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	devices, err := core.GetEngine().GetPresence().GetClientOnlineDevices(clientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Presence: failed to marshal response %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := onlineDevicesResponse{
		OnlineDeviceIDs: devices,
	}

	data, err := json.Marshal(&response)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Presence: failed to marshal response %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}

// GetClientDevicesPresences - Get all client devices presence
// GET /presence/:clientID
func GetClientDevicesPresences(context *gin.Context) {
	writer := context.Writer
	request := context.Request

	tokenString := request.Header.Get("Authorization")

	// Token must be present
	if tokenString == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	identity, isOK := auth.VerifyToken(tokenString)

	// Invalid token
	if !isOK {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// If not super admin, then ther client doesn't have enough permissions
	if !identity.IsSuperAdmin() {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	clientID := context.Params.ByName("clientID")

	if clientID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	devices, err := core.GetEngine().GetPresence().GetClientDevicesPresences(clientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Presence: failed to fetch client online devices %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := lastDeviceStatusResponse{
		Devices: devices,
	}

	data, err := json.Marshal(&response)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Presence: failed to marshal response %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}
