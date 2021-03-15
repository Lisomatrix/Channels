package auth

import "net/http"

// SuperAdminRole - Super Admin Role
var SuperAdminRole = "SuperAdmin"

// AdminRole - Admin Role
var AdminRole = "Admin"

// ClientRole - Client Role
var ClientRole = "Client"

// Identity - Data in JWT token
type Identity struct {
	Role     string
	AppID    string
	ClientID string
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
