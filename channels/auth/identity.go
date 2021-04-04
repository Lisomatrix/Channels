// Provides the structs and function to verify and extract a token
// and check if the token can use the resources he wants
package auth

import (
	"net/http"
)

// All available roles
const (
	// SuperAdminRole - Super Admin Role
	SuperAdminRole = "SuperAdmin"

	// AdminRole - Admin Role
	AdminRole = "Admin"

	// ClientRole - Client Role
	ClientRole = "Client"
)

// Identity - Data in JWT token
type Identity struct {
	Role     string
	AppID    string
	ClientID string
}

func GetTokenAndVerify(request *http.Request, role string) (*Identity, bool) {
	tokenString := request.Header.Get("Authorization")

	// Token must be present
	if tokenString == "" {
		return nil, false
	}

	if identity, isOK := VerifyToken(tokenString); !isOK {
		return nil, false
	} else if identity.Role != role {
		return &identity, false
	} else {
		return &identity, false
	}
}

// GetAuthData - Get AppID and Authorization
func GetAuthData(request *http.Request) (string, string, bool) {
	// Get AppID and Token
	appID := request.Header.Get("AppID")
	token := request.Header.Get("Authorization")

	if token == "" || appID == "" {
		return "", "", false
	}

	return token, appID, true
}

// CanUseAppID - Check if user can use given AppID
func (identity *Identity) CanUseAppID(appID string) bool {

	if identity.IsSuperAdmin() {
		return true
	}

	if identity.AppID == appID {
		return true
	}

	return false
}

// IsAdminKind - Check if is Super Admin or Admin
func (identity *Identity) IsAdminKind() bool {
	return identity.IsSuperAdmin() || identity.IsAdmin()
}

// IsSuperAdmin - Check if client is a Super Admin
func (identity *Identity) IsSuperAdmin() bool {
	return identity.Role == SuperAdminRole
}

// IsAdmin - Check if client is admin only
func (identity *Identity) IsAdmin() bool {
	return identity.Role == AdminRole
}

// IsClient - Check if client is not admin
func (identity *Identity) IsClient() bool {
	return identity.Role == ClientRole
}
