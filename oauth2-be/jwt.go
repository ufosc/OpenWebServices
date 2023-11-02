package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Claims are the JWT claims extracted from a cookie.
type Claims struct {
	ID    string `json:"id"`
	PHash string `json:"phash"`
}

// NewJWT creates a jwt token string for the given user, with the given secret
// string. Returns empty string on error.
func NewJWT(secret string, user UserModel) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"phash": user.Password,
		"exp":   jwt.NewNumericDate(time.Now().Add(20 * time.Minute)).Unix(),
	})

	str, err := token.SignedString(secret)
	if err != nil {
		return ""
	}

	return str
}

// ValidateJWT checks if the JWT string literal is valid. It does not check if
// the user password has been changed.
func ValidateJWT(literal string, config Config) (Claims, bool) {
	token, err := jwt.Parse(literal, func(token *jwt.Token) (interface{}, error) {

		// Signing methods must be the same.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Return the signing secret.
		return config.SECRET, nil
	})

	if err != nil {
		return Claims{}, false
	}

	// TODO: does this check if token has expired?
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return Claims{}, false
	}

	// Cast claims.
	id, ok := claims["id"].(string)
	if !ok {
		return Claims{}, false
	}

	phash, ok := claims["phash"].(string)
	if !ok {
		return Claims{}, false
	}

	return Claims{id, phash}, true
}
