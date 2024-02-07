package authmw

import (
	"github.com/gin-gonic/gin"
	"github.com/ufosc/OpenWebServices/pkg/authdb"
	"time"
)

// A returns a middleware function that authorizes a route using
// an assertion parameter (access token).
func A(db authdb.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		assertion := c.DefaultQuery("assertion", "")
		if assertion == "" {
			setError(c, ErrInvalid, "missing assertion")
			return
		}

		// Verify access token exists.
		tkExists, err := db.Tokens().FindAccessByID(assertion)
		if err != nil {
			setError(c, ErrToken, "access token expired/not found")
			return
		}

		// Key must be issued from dashboard (frontend).
		if tkExists.ClientID != "0" {
			setError(c, ErrToken, "access token must be from dashboard")
			return
		}

		// Ensure key is not expired.
		if (tkExists.CreatedAt + tkExists.TTL) < time.Now().Unix() {
			db.Tokens().DeleteAccessByID(assertion)
			setError(c, ErrToken, "Access token expired / not found")
			return
		}

		// Verify associated user exists.
		userExists, err := db.Users().FindByID(tkExists.UserID)
		if err != nil {
			db.Tokens().DeleteAccessByID(assertion)
			setError(c, ErrToken, "User not found")
			return
		}

		c.Set("user", userExists)
		c.Next()
	}
}
