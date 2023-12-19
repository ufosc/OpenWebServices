package common

import (
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"regexp"
	pwd "github.com/wagslane/go-password-validator"
)

// ValidateEmail checks whether email is a valid email address.
func ValidateEmail(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}
	return true
}

// ValidatePassword checks whether password is strong enough.
func ValidatePassword(password string) error {
	return pwd.Validate(password, 60)
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
