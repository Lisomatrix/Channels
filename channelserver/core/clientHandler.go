package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"lisomatrix.pt/channelserver/channelserver/auth"
)

type createClientRequest struct {
	ClientID string `json:"clientID"`
	Username string `json:"username"`
	Extra    string `json:"extra"`
}

type updateClientRequest struct {
	Username string `json:"username"`
	Extra    string `json:"extra"`
}

type getClientsResponse struct {
	Clients []*Client `json:"clients"`
}

// CreateClientHandler - Create new client
// POST /client
func CreateClientHandler(context *gin.Context) {

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
	var createClientRequest createClientRequest

	err = json.Unmarshal(body, &createClientRequest)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	existingClient, err := GetClient(appID, createClientRequest.ClientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Create Client failed to check if client already exists %v\n", err)
	} else if existingClient != nil {
		writer.WriteHeader(http.StatusConflict)
		return
	}

	if isOK, err := CreateClient(appID, createClientRequest.ClientID, createClientRequest.Username, createClientRequest.Extra); err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Create Client failed %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
	} else if isOK {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusConflict)
	}
}

// DeleteClientHandler - Delete a client
// DELETE /client/:clientID
func DeleteClientHandler(context *gin.Context) {

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

	clientID := context.Params.ByName("clientID")

	if clientID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if isOK, err := DeleteClient(appID, clientID); err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Delete Client: failed to delete client %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
	} else if isOK {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

// UpdateClientHandler - Update a client
// PUT /client/:clientID
func UpdateClientHandler(context *gin.Context) {

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
	var updateClientRequest updateClientRequest

	err = json.Unmarshal(body, &updateClientRequest)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	clientID := context.Params.ByName("clientID")

	if clientID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, err := GetEngine().GetClientRepository().ExistsAppClient(appID, clientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Update Client: failed to check client existence %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	err = GetEngine().GetClientRepository().UpdateClient(clientID, updateClientRequest.Username, updateClientRequest.Extra)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Update Client failed %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

// GetClients - Get app or all the clients
// GET /client
func GetClients(context *gin.Context) {

	request := context.Request
	writer := context.Writer

	// Check for required headers
	// Get AppID and Token
	appID := request.Header.Get("AppID")
	token := request.Header.Get("Authorization")

	if token == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(token)

	// If not valid return
	if !isOK {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	var getClientsResponse getClientsResponse

	if appID != "" && identity.CanUseAppID(appID) {
		clients, err := GetEngine().GetClientRepository().GetAppClients(appID)

		if err != nil {
			fmt.Fprintf(os.Stderr, "HTTP Get app clients failed %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		getClientsResponse.Clients = clients
	} else if appID == "" && identity.IsSuperAdmin() {
		clients, err := GetEngine().GetClientRepository().GetAllClients()

		if err != nil {
			fmt.Fprintf(os.Stderr, "HTTP Get all clients failed %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		getClientsResponse.Clients = clients
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := json.Marshal(getClientsResponse)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get all clients: Marshal failed %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Write(data)
	writer.WriteHeader(http.StatusOK)
}

// GetClientHandler - Get client info
// GET /client/:clientID
func GetClientHandler(context *gin.Context) {

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

	clientID := context.Params.ByName("clientID")

	if clientID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	client, err := GetClient(appID, clientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get Client failed %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(client)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Get Client: failed to marshal response %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Write(data)
	writer.WriteHeader(http.StatusOK)
}
