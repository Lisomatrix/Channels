package core

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/lisomatrix/channels/channels/auth"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type createDeviceRequest struct {
	DeviceID string `json:"deviceID"`
	Token    string `json:"token"`
}

type updateDeviceRequest struct {
	Token string `json:"token"`
}

// CreateDevice - Add new device for notifications
// POST /device
func CreateDevice(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Check for required headers
	token, appID, isOK := auth.GetAuthData(request)

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if is admin, and validate token
	identity, isOK := auth.VerifyToken(token)

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
	var createDeviceRequest createDeviceRequest

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(body, &createDeviceRequest)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	exists := GetEngine().GetCacheStorage().CheckDeviceExistence(identity.ClientID, createDeviceRequest.DeviceID)

	if !exists {

		device, err := GetEngine().GetDeviceRepository().GetDevice(createDeviceRequest.DeviceID)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "HTTP Create Device: failed to check device existence %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if device != nil {
			writer.WriteHeader(http.StatusConflict)

			GetEngine().GetCacheStorage().AddDevice(device.ClientID, device)

			return
		}

	} else {
		writer.WriteHeader(http.StatusConflict)
		return
	}

	device := &Device{
		ID:       createDeviceRequest.DeviceID,
		Token:    createDeviceRequest.Token,
		ClientID: identity.ClientID,
	}

	err = GetEngine().GetDeviceRepository().CreateDevice(device.ID, device.Token, device.ClientID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "HTTP Create Device: failed to create device %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	GetEngine().GetCacheStorage().AddDevice(identity.ClientID, device)

	writer.WriteHeader(http.StatusOK)
}

// RemoveDevice - Remove device for notifications
// DELETE /device/:deviceID
func RemoveDevice(context *gin.Context) {
	request := context.Request
	writer := context.Writer

	// Check for required headers
	token, appID, isOK := auth.GetAuthData(request)

	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if is admin, and validate token
	identity, isOK := auth.VerifyToken(token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(appID) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get Url param
	deviceID := context.Param("deviceID")

	if deviceID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Fetch device from database
	device, err := GetEngine().GetDeviceRepository().GetDevice(deviceID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "HTTP Delete Device: failed to check device existence %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if device == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// Check that device belongs to this client
	if device.ClientID != identity.ClientID {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Delete the device
	err = GetEngine().GetDeviceRepository().DeleteDevice(deviceID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "HTTP DELETE Device: failed to delete device %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Remove cache
	GetEngine().GetCacheStorage().RemoveDevice(identity.ClientID, deviceID)

	writer.WriteHeader(http.StatusOK)
}
