package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AuthenticateUser(db *Database, config Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Retrieve auth cookie.
		cookie, err := c.Cookie("ows-jwt")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrNoCookie,
			})
			return
		}

		// Extract JWT claims. Will fail if expired.
		claims, ok := ValidateJWT(cookie, config)
		if !ok {
			c.SetCookie("ows-jwt", "", 0, "/", "localhost", false, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrNoCookie,
			})
			return
		}

		// Check if user exists.
		userExists, err := db.ReadUser(claims.ID)
		if err != nil {
			c.SetCookie("ows-jwt", "", 0, "/", "localhost", false, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrUserNotFound,
			})
			return
		}

		// Ensure user password did not change.
		if userExists.Password != claims.PHash {
			c.SetCookie("ows-jwt", "", 0, "/", "localhost", false, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrPasswordChanged,
			})
			return
		}

		// Ensure LastVerified is less than 3 months ago.
		lastVerified := time.Unix(userExists.LastVerified, 0)
		if time.Since(lastVerified).Hours() > 2160 {
			c.SetCookie("ows-jwt", "", 0, "/", "localhost", false, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrVerificationRequired,
			})
			return
		}

		c.Set("user", userExists)
	}
}

func AuthenticateClient(db *Database) gin.HandlerFunc {
	// TODO: verify client.Owner still exists.
	return func(c *gin.Context) {
	}
}

func AuthenticateToken(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
