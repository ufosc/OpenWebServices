package main

import (
	"net/mail"
	"strings"
)

const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numeric = "0123456789"
const special = "`~!@#$%^&*()_-=+[{]}\\|:;\"',<.>/"

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

func VerifyPassword(hash, password string) bool {
	return true
}
