package main

import (
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"strings"
)

const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const special = "`~!@#$%^&*()_-=+[{]}\\|:;\"',<.>/"
const numeric = "0123456789"

// ValidateEmail checks whether email is a valid email address ending
// with the @ufl.edu domain.
func ValidateEmail(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	if parts[1] != "ufl.edu" {
		return false
	}

	return true
}

// ValidatePassword checks whether password is at least 12 chars, contains
// both digits and characters, and has at least 1 special character.
func ValidatePassword(password string) string {
	// Must be a minimum of 12 characters.
	if len(password) < 12 {
		return "password must be at least 12 characters"
	}

	// Must be alphanumeric (both).
	if !strings.ContainsAny(password, alpha) {
		return "password must contain at least one letter"
	}

	if !strings.ContainsAny(password, numeric) {
		return "password must contain at least one digit"
	}

	// Must have 1 special character.
	if !strings.ContainsAny(password, special) {
		return "password must contain at least one special character"
	}

	return ""
}

// VerifyPassword verifies that a given password matches the hash.
func VerifyPassword(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

func validateTokenScope(scope []string) bool {
	// Token-based clients only have access to 'public' scope.
	if len(scope) > 1 {
		return false
	}

	if scope[0] != "public" {
		return false
	}

	return true
}

func validateCodeScope(scope []string) bool {
	if len(scope) > 2 {
		return false
	}

	// Must be either 'public' or 'email'
	for _, v := range scope {
		if v != "public" && v != "email" {
			return false
		}
	}

	return true
}

// ValidateScope checks whether the scope string is valid for the given response
// type.
func ValidateScope(resType string, scope []string) bool {
	if len(scope) == 0 {
		return false
	}

	if resType == "token" {
		return validateTokenScope(scope)
	}

	if resType == "code" {
		return validateCodeScope(scope)
	}

	return false
}
