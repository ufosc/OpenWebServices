package common

import (
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"regexp"
	"strings"
)

const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const special = "`~!@#$%^&*()_-=+[{]}\\|;\"',<.>/"
const numeric = "0123456789"

// ValidateEmail checks whether email is a valid email address.
func ValidateEmail(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
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

	// TODO: this is checking whether the password contains valid
	// characters, not whether all characters are valid.

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

// ValidateTokenScope validates a scope slice for a client that uses the 'token'
// authentication response type.
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

// ValidateCodeScope validates a scope slice for a client that uses the 'code'
// authentication response type.
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

var redirectURIRegex = regexp.MustCompile(`^https:\/\/([a-z0-9]|\.|\-|\_)*(\/[a-z0-9A-Z:@?=&%.\-_$+]+)*$`)

// ValidateRedirectURI validates a client redirect URI.
func ValidateRedirectURI(uri string) bool {
	return !redirectURIRegex.MatchString(uri)
}
