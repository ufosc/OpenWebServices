package authmw

import (
	"github.com/gin-gonic/gin"
	"github.com/ufosc/OpenWebServices/pkg/authdb"
	"github.com/ufosc/OpenWebServices/pkg/common"
	"net/http"
	"strings"
	"time"
)

// AuthenticateUser is a middleware that verifies a user JWT in the
// assertion URL parameter.
func AuthenticateUser(secret string, db authdb.Database,
	realms ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAssertion := c.DefaultQuery("assertion", "")
		if userAssertion == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "expected assertion URL parameter",
			})
			return
		}

		// Validate JWT.
		claims, ok := common.ValidateJWT(userAssertion, secret)
		if !ok || claims.Type != "user" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired assertion",
			})
			return
		}

		// Check if subject exists.
		userExists, err := db.Users().FindByID(claims.Sub)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user not found",
			})
			return
		}

		// Ensure password hash not changed.
		if userExists.Password != claims.PKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user password has changed",
			})
			return
		}

		// Ensure user has all required realms.
		hasRealms := map[string]bool{}
		for _, realm := range userExists.Realms {
			hasRealms[realm] = true
		}

		for _, realm := range realms {
			if !hasRealms[realm] {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "insufficient user permission",
				})
				return
			}
		}

		c.Set("user", userExists)
		c.Next()
	}
}

// AuthenticateClient is a middleware that authenticates a client JWT in
// the client_assertion URL parameter.
func AuthenticateClient(secret string, db authdb.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientAssertion := c.DefaultQuery("client_assertion", "")
		if clientAssertion == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "expected assertion URL parameter",
			})
			return
		}

		// Validate JWT.
		claims, ok := common.ValidateJWT(clientAssertion, secret)
		if !ok || claims.Type != "client" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired client_assertion",
			})
			return
		}

		// Ensure client exists.
		clientExists, err := db.Clients().FindByID(claims.Sub)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "client not found",
			})
			return
		}

		// Ensure client owner still exists.
		_, err = db.Users().FindByID(clientExists.Owner)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user owner associated with client no longer exists",
			})
			return
		}

		// Ensure client is not expired.
		if clientExists.CreatedAt+clientExists.TTL < time.Now().Unix() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "client lease has expired",
			})
			return
		}

		// Ensure key has not changed.
		if clientExists.Key != claims.PKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid client key",
			})
			return
		}

		c.Set("client", clientExists)
		c.Next()
	}
}

// AuthenticateBearer is a middleware the verifies an access bearer token.
func AuthenticateBearer(secret string, db authdb.Database,
	scope ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scopeStr := ""
		for _, val := range scope {
			scopeStr += val
		}

		tkStr := strings.Split(c.GetHeader("Authorization"), " ")
		if len(tkStr) < 2 {
			c.Header("WWW-Authenticate", "Bearer scope=\""+scopeStr+
				"\" error=\"invalid_request\"")

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "expected authorization header",
			})
			return
		}

		if tkStr[0] != "Bearer" {
			c.Header("WWW-Authenticate", "Bearer scope=\""+scopeStr+
				"\" error=\"invalid_request\"")

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization scheme must be 'bearer'",
			})
			return
		}

		// Check if token is JWT.
		if claims, ok := common.ValidateJWT(tkStr[1], secret); ok {
			if claims.Type != "user" {
				c.Header("WWW-Authenticate", "Bearer scope=\""+scopeStr+
					"\" error=\"invalid_request\"")

				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "jwt must be of type 'user'",
				})

				return
			}

			// Check if subject exists.
			userExists, err := db.Users().FindByID(claims.Sub)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "user not found",
				})
				return
			}

			// Ensure password hash not changed.
			if userExists.Password != claims.PKey {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "user password has changed",
				})
				return
			}

			c.Set("user", userExists)
			c.Set("client", authdb.ClientModel{
				Scope: []string{"email", "public"},
			})

			return
		}

		// Verify key exists.
		tkExists, err := db.Tokens().FindAccessByID(tkStr[1])
		if err != nil {
			c.Header("WWW-Authenticate", "Bearer scope=\""+scopeStr+
				"\" error=\"invalid_token\"")

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Access token expired or could not be found",
			})
			return
		}

		// Ensure key is not expired.
		if (tkExists.CreatedAt + tkExists.TTL) > time.Now().Unix() {
			c.Header("WWW-Authenticate", "Bearer scope=\""+scopeStr+
				"\" error=\"invalid_token\"")

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Access token expired or could not be found",
			})
			return
		}

		// Verify associated user exists.
		userExists, err := db.Users().FindByID(tkExists.UserID)
		if err != nil {
			c.Header("WWW-Authenticate", "Bearer scope=\""+scopeStr+
				"\" error=\"invalid_token\"")

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "The user associated with this token could not be found",
			})
			return
		}

		// Verify associated client exists.
		clientExists, err := db.Clients().FindByID(tkExists.ClientID)
		if err != nil {
			c.Header("WWW-Authenticate", "Bearer scope=\""+scopeStr+
				"\" error=\"invalid_token\"")

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "The client associated with this token could not be found",
			})
			return
		}

		// Verify token is authorized for this route.
		haveScope := map[string]bool{}
		for _, value := range clientExists.Scope {
			haveScope[value] = true
		}

		for _, value := range scope {
			if !haveScope[value] {
				c.Header("WWW-Authenticate", "Bearer scope=\""+scopeStr+
					"\" error=\"insufficient_scope\"")

				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Insufficient client permission",
				})

				return
			}
		}

		// Write client, user, token to context.
		c.Set("user", userExists)
		c.Set("client", clientExists)
		c.Next()
	}
}
