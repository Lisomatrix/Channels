package auth

import (
	"fmt"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

var jwtSecret = []byte("")
var issuer = ""

func SetSecret(secret string) {
	jwtSecret = []byte(secret)
}

func SetIssuer(issuer string) {
	issuer = issuer
}

// CreateToken - Create token with given info
func CreateToken(clientID string, role string, appID string, expire *jwt.NumericDate) (string, error) {

	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: jwtSecret}, (&jose.SignerOptions{}).WithType("JWT"))

	c := &jwt.Claims{
		Subject: clientID,
		Issuer:  issuer,
		Expiry: expire,
	}

	tokenData := Identity{
		AppID: appID,
		ClientID: clientID,
		Role: role,
	}

	raw, err := jwt.Signed(sig).Claims(c).Claims(tokenData).CompactSerialize()

	if err != nil {
		fmt.Printf("Error creating jwt %s \n", err)
		return "", err
	}

	return raw, nil
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

	token, err := jwt.ParseSigned(tokenString)

	if err != nil {
		return tokenData, false
	}

	out := Identity{}

	if err := token.Claims(jwtSecret, &out); err != nil {
		return out, false
	}


	return out, true
}
