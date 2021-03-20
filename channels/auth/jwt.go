package auth

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("")

func SetSecret(secret string) {
	jwtSecret = []byte(secret)
}

// AuthenticateAdmin - Check if token is valid and is admin kind
func AuthenticateAdmin(tokenString string) (*Identity, bool) {

	// Verify token
	identity, isOK := VerifyToken(tokenString)

	// If token is invalid
	if !isOK {
		return nil, false
	}

	// If is not admin type
	if !identity.IsAdminKind() {
		return nil, false
	}

	return &identity, true
}

// VerifyToken - Check token validity and payload
func VerifyToken(tokenString string) (Identity, bool) {

	var tokenData Identity
	isOK := false

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})

	if err != nil {
		return tokenData, isOK
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get data from jwt
		role := claims["Role"].(string)
		clientID := claims["ClientID"].(string)

		var appID string

		if claims["AppID"] == nil {
			appID = ""
		} else {
			appID = claims["AppID"].(string)
		}

		// Check if there is a role and clientID
		if role == "" || clientID == "" {
			isOK = false

			// If so, check if user is a SuperAdmin
		} else if role == SuperAdminRole {
			tokenData = Identity{
				Role:     role,
				ClientID: clientID,
				AppID: appID,
			}
			isOK = true
			// Otherwise check if user is Admin or Client, and check for AppID presence
		} else if (role == AdminRole || role == ClientRole) && appID != "" {
			tokenData = Identity{
				Role:     role,
				ClientID: clientID,
				AppID:    appID,
			}
			isOK = true
		} else {
			// If the requirements were not meet, then this token is not valid
			isOK = false
		}

	}

	return tokenData, isOK
}
