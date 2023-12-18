package common

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ufosc/OpenWebServices/pkg/authdb"
	"time"
)

// Claims are the JWT claims.
type Claims struct {
	Sub  string `json:"sub"`
	Type string `json:"type"`
	PKey string `json:"pkey"`
	Iss  string `json:"iss"`
	Aud  string `json:"aud"`
	Exp  int64  `json:"exp"`
}

// NewUserJWT creates a jwt token string for the given user, with the given secret
// string. Returns empty string on error.
func NewUserJWT(secret string, user authdb.UserModel) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"type": "user",
		"pkey": user.Password,
		"iss":  "ufosc-jwt",
		"aud":  "https://api.ufosc.org/auth/authorize",
		"exp":  time.Now().Add(20 * time.Minute).Unix(),
	})

	str, err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}

	return str
}

// NewClientJWT creates a jwt token string for the given client, with the given
// secret string. Returns empty string on error.
func NewClientJWT(secret string, client authdb.ClientModel) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  client.ID,
		"type": "client",
		"pkey": client.Key,
		"iss":  "ufosc-jwt",
		"aud":  "https://api.ufosc.org/auth/authorize",
		"exp":  time.Now().Add(20 * time.Minute).Unix(),
	})

	str, err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}

	return str
}

// ValidateJWT checks if the JWT string literal is valid. It does not check if
// the user password has been changed.
func ValidateJWT(literal, secret string) (Claims, bool) {
	token, err := jwt.Parse(literal, func(token *jwt.Token) (interface{}, error) {

		// Signing methods must be the same.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the signing secret.
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println(err)
		return Claims{}, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return Claims{}, false
	}

	// Cast claims.
	sub, ok := claims["sub"].(string)
	if !ok {
		return Claims{}, false
	}

	ttype, ok := claims["type"].(string)
	if !ok {
		return Claims{}, false
	}

	pkey, ok := claims["pkey"].(string)
	if !ok {
		return Claims{}, false
	}

	iss, ok := claims["iss"].(string)
	if !ok {
		return Claims{}, false
	}

	aud, ok := claims["aud"].(string)
	if !ok {
		return Claims{}, false
	}

	// Validate claims.
	if ttype != "client" && ttype != "user" {
		return Claims{}, false
	}

	if iss != "ufosc-jwt" || aud != "https://api.ufosc.org/auth/authorize" {
		return Claims{}, false
	}

	return Claims{sub, ttype, pkey, iss, aud, 0}, true
}
