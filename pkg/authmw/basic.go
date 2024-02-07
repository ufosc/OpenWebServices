package authmw

import (
	b64 "encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/ufosc/OpenWebServices/pkg/authdb"
	"github.com/ufosc/OpenWebServices/pkg/common"
	"strings"
	"time"
)

// B returns a middleware that implements the basic authorization
// scheme for clients.
func B(db authdb.Database) gin.HandlerFunc {
	// See: https://datatracker.ietf.org/doc/html/rfc6749#section-6
	return func(c *gin.Context) {
		tkStr := strings.Split(c.GetHeader("Authorization"), " ")
		if len(tkStr) != 2 {
			setError(c, ErrInvalid, "expected Authorization header")
			return
		}

		if tkStr[0] != "Basic" {
			setError(c, ErrInvalid, "auth scheme must be Basic")
			return
		}

		// Decode token.
		sDec, err := b64.StdEncoding.DecodeString(tkStr[1])
		if err != nil {
			setError(c, ErrInvalid, "invalid authorization token")
			return
		}

		parts := strings.Split(string(sDec), ":")
		if len(parts) != 2 {
			setError(c, ErrInvalid, "invalid authorization token")
			return
		}

		// Verify client exists.
		clientExists, err := db.Clients().FindByID(parts[0])
		if err != nil {
			setError(c, ErrToken, "client not found")
			return
		}

		// Client must be secure (code type).
		if clientExists.ResponseType != "code" {
			setError(c, ErrToken, "client is insecure")
			return
		}

		// Client must not be expired.
		if (clientExists.CreatedAt + clientExists.TTL) <
			time.Now().Unix() {
			db.Clients().DeleteByID(clientExists.ID)
			setError(c, ErrToken, "client has expired")
			return
		}

		// Client owner must still exist.
		if _, err := db.Users().FindByID(clientExists.Owner); err != nil {
			db.Clients().DeleteByID(clientExists.ID)
			setError(c, ErrToken, "client owner account no longer exists")
			return
		}

		// Verify password.
		if !common.VerifyPassword(clientExists.Key, parts[1]) {
			setError(c, ErrToken, "incorrect key")
			return
		}

		c.Set("client", clientExists)
		c.Next()
	}
}
