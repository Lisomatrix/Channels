package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Channels/Channels/auth"
	"github.com/gin-gonic/gin"
)

type createAppRequest struct {
	AppID string `json:"appID"`
	Name  string `json:"name"`
}

// CreateAppHandler - Create a new App
func CreateAppHandler(context *gin.Context) {
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

	body, err := ioutil.ReadAll(request.Body)

	// If no body, then request is invalid
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var createRequest createAppRequest

	// Parse body
	err = json.Unmarshal(body, &createRequest)

	// If there is and error, then body is invalid
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}
