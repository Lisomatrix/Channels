package auth

import "testing"

func TestCanUserAppID(t *testing.T) {
	// clientID: 321
	// appID: 123

	superAdmin := Identity{
		Role:     SuperAdminRole,
		AppID:    "123",
		ClientID: "321",
	}

	if !superAdmin.CanUseAppID("123") {
		t.Errorf("Failed to verify if super admin can use given AppID")
	}

	// Super admins can use all AppIDs
	if !superAdmin.CanUseAppID("321") {
		t.Errorf("Failed to verify if super admin can use given AppID")
	}

	if !superAdmin.IsAdminKind() {
		t.Errorf("Failed to verify if super admin is admin kind")
	}

	if !superAdmin.IsSuperAdmin() {
		t.Errorf("Failed to verify if super admin is super admin")
	}

	if superAdmin.IsAdmin() {
		t.Errorf("Failed to verify if super admin is admin")
	}

	if superAdmin.IsClient() {
		t.Errorf("Failed to verify if super admin is client")
	}

	admin := Identity{
		Role:     AdminRole,
		AppID:    "123",
		ClientID: "321",
	}

	if !admin.CanUseAppID("123") {
		t.Errorf("Failed to verify if admin can use given AppID")
	}

	if admin.CanUseAppID("321") {
		t.Errorf("Failed to verify if admin can use given AppID, he shouldn't have access")
	}

	if !admin.IsAdminKind() {
		t.Errorf("Failed to verify if admin is admin kind")
	}

	if admin.IsSuperAdmin() {
		t.Errorf("Failed to verify if admin is super admin")
	}

	if !admin.IsAdmin() {
		t.Errorf("Failed to verify if admin is admin")
	}

	if admin.IsClient() {
		t.Errorf("Failed to verify if admin is client")
	}

	client := Identity{
		Role:     ClientRole,
		AppID:    "123",
		ClientID: "321",
	}

	if !client.CanUseAppID("123") {
		t.Errorf("Failed to verify if client can use given AppID, he shouldn't have access")
	}

	if client.CanUseAppID("321") {
		t.Errorf("Failed to verify if client can use given AppID, he shouldn't have access")
	}

	if client.IsAdminKind() {
		t.Errorf("Failed to verify if client is admin kind")
	}

	if client.IsSuperAdmin() {
		t.Errorf("Failed to verify if client is super admin")
	}

	if client.IsAdmin() {
		t.Errorf("Failed to verify if client is admin")
	}

	if !client.IsClient() {
		t.Errorf("Failed to verify if client is client")
	}
}

func TestVerifyAdmin(t *testing.T) {
	// clientID: 321
	// appID: 123
	SetSecret("123")
	superAdminToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJSb2xlIjoiU3VwZXJBZG1pbiIsIkNsaWVudElEIjoiMzIxIiwiQXBwSUQiOiIxMjMifQ.k5Yf1E_oqCYUJndDM1pJ2aRISpxGGNMHIMt4EGgOQKw"
	adminToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJSb2xlIjoiQWRtaW4iLCJDbGllbnRJRCI6IjMyMSIsIkFwcElEIjoiMTIzIn0.pEvQg3zVCw0zuIM0QCcDsHim9CgMRh5zc5YU3bDaOTg"
	clientToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJSb2xlIjoiQ2xpZW50IiwiQ2xpZW50SUQiOiIzMjEiLCJBcHBJRCI6IjEyMyJ9.nfDURuXKPU_GVK0TgQMzXSrfsdyiFO6jgn-wOPDDjF4"

	_, isOK := AuthenticateAdmin(superAdminToken)

	if !isOK {
		t.Errorf("Failed to verify super admin token")
	}

	_, isOK = AuthenticateAdmin(adminToken)

	if !isOK {
		t.Errorf("Failed to verify admin token")
	}

	_, isOK = AuthenticateAdmin(clientToken)

	if isOK {
		t.Errorf("Client token is validated as admin token!")
	}
}

func TestVerifyToken(t *testing.T) {
	// clientID: 321
	// appID: 123
	// role: SuperAdmin
	SetSecret("123")
	superAdminToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJSb2xlIjoiU3VwZXJBZG1pbiIsIkNsaWVudElEIjoiMzIxIiwiQXBwSUQiOiIxMjMifQ.k5Yf1E_oqCYUJndDM1pJ2aRISpxGGNMHIMt4EGgOQKw"

	identity, isOK := VerifyToken(superAdminToken)

	if !isOK {
		t.Errorf("Failed to verify valid super admin token")
		return
	}

	if identity.AppID != "123" {
		t.Errorf("Admin token appID doesn't match")
	}

	if identity.ClientID != "321" {
		t.Errorf("Admin token clientID doesn't match")
	}

	if identity.Role != SuperAdminRole {
		t.Errorf("Admin token role doesn't match")
	}
}

// TestCreateToken - Test create token and verify
func TestCreateToken(t *testing.T) {
	clientID := "123"
	role := ClientRole
	appID := "123"

	SetSecret("123")

	token, err := CreateToken(clientID, role, appID, nil)

	if err != nil {
		t.Errorf("Failed to create token! \n")
		return
	}

	identity, isOK := VerifyToken(token)

	if !isOK {
		t.Errorf("Failed to verify created token! \n")
		return
	}

	if identity.AppID != appID {
		t.Errorf("Token AppID doesn't match the given one")
	}

	if identity.ClientID != clientID {
		t.Errorf("Token ClientID doesn't match the given one")
	}

	if identity.Role != role {
		t.Errorf("Token Role doesn't match the given one")
	}
}
