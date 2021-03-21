// This package hold all the delivery logic, connection tracking and pub sub
package core

import (
	jsoniter "github.com/json-iterator/go"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/lisomatrix/channels/channels/auth"

	"github.com/gin-gonic/gin"
)

type createAppRequest struct {
	AppID string
	Name  string
}

type updateAppRequest struct {
	Name string
}

type getAppsResponse struct {
	Apps []*App
}

// CreateApp - Create a new app
// POST /app
func CreateApp(context *gin.Context) {

	request := context.Request
	writer := context.Writer

	// Get AppID and Token
	token := request.Header.Get("Authorization")

	if token == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(token)

	// If not valid return
	if !isOK || !identity.IsSuperAdmin() {
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
	var createAppRequest createAppRequest

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(body, &createAppRequest)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if app already exists in cache
	app := GetEngine().GetCacheStorage().GetApp(createAppRequest.AppID)

	// If not found check from database
	if app == nil {
		exists, err := GetEngine().GetAppRepository().AppExists(createAppRequest.AppID)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "HTTP Create App: failed to check app existence %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if exists {
			writer.WriteHeader(http.StatusConflict)
			return
		}
	}

	// Store app in the database
	err = GetEngine().GetAppRepository().CreateApp(createAppRequest.AppID, createAppRequest.Name)

	// Store app in cache
	GetEngine().GetCacheStorage().StoreApp(createAppRequest.AppID, createAppRequest.Name)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "HTTP Create App: failed to create app %v\n", err)
		writer.WriteHeader(http.StatusConflict)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

// DeleteApp - Delete a app
// DELETE /app/:appID
func DeleteApp(context *gin.Context) {

	request := context.Request
	writer := context.Writer

	// Get AppID and Token
	token := request.Header.Get("Authorization")

	if token == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(token)

	// If not valid return
	if !isOK || !identity.IsSuperAdmin() {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	appID := context.Params.ByName("appID")

	if appID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err := GetEngine().GetAppRepository().DeleteApp(appID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "HTTP Delete App: failed to delete app %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	GetEngine().GetCacheStorage().RemoveApp(appID)

	writer.WriteHeader(http.StatusOK)
}

// UpdateApp - Update app
// PUT /app/:appID
func UpdateApp(context *gin.Context) {

	request := context.Request
	writer := context.Writer

	// Get AppID and Token
	token := request.Header.Get("Authorization")

	if token == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(token)

	// If not valid return
	if !isOK {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	appID := context.Params.ByName("appID")

	if appID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// If client is just admin then check it is their app
	if !identity.CanUseAppID(appID) {
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
	var updateAppRequest updateAppRequest

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(body, &updateAppRequest)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if app already exists in cache
	app := GetEngine().GetCacheStorage().GetApp(appID)

	// If app not found in cache
	if app == nil {
		// Check its exitence in the database
		exists, err := GetEngine().GetAppRepository().AppExists(appID)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "HTTP Update App: failed to check app existence %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// If not found then the request is invalid
		if !exists {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
	}

	// Update app
	err = GetEngine().GetAppRepository().UpdateApp(appID, updateAppRequest.Name)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "HTTP Update App: failed to update app %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update cache entry
	GetEngine().GetCacheStorage().StoreApp(appID, updateAppRequest.Name)

	writer.WriteHeader(http.StatusOK)
}

// GetApps - Get all apps
// Get /app
func GetApps(context *gin.Context) {

	request := context.Request
	writer := context.Writer

	// Get AppID and Token
	token := request.Header.Get("Authorization")

	if token == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(token)

	// If not valid return
	if !isOK || !identity.IsSuperAdmin() {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get alls apps
	apps, err := GetEngine().GetAppRepository().GetApps()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "HTTP Get Apps: failed to get apps %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Since we fetched the apps, let use the opportunity to cache them
	go func() {
		for _, app := range apps {
			GetEngine().GetCacheStorage().StoreApp(app.AppID, app.Name)
		}
	}()

	response := getAppsResponse{Apps: apps}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(response)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "HTTP Get Apps: failed to marshal response %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(data)
}
