package authmw

import (
	"github.com/gin-gonic/gin"
	"github.com/ufosc/OpenWebServices/pkg/authdb"
	"net/http"
	"strings"
	"time"
)

// Config defines the scope and realm requirements for a
// route authentication middleware.
type Config struct {
	Scope  []string
	Realms []string
}

// WWW-Authenticate response header errors.
// See: https://datatracker.ietf.org/doc/html/rfc6750#section-3
const (
	ErrInvalid = "invalid_request"
	ErrToken   = "invalid_token"
	ErrScope   = "insufficient_scope"
)

func setError(c *gin.Context, code, desc string) {
	scopes, _ := c.Get("header-scopes")
	realms, _ := c.Get("header-realms")

	scopesStr, _ := scopes.(string)
	realmsStr, _ := realms.(string)

	c.Header("WWW-Authenticate", "Bearer scope=\""+scopesStr+
		"\", realms=\""+realmsStr+"\", error=\""+code+
		"\", error_description=\""+desc+"\"")

	httpCode := http.StatusUnauthorized
	if code == ErrInvalid {
		httpCode = http.StatusBadRequest
	}

	c.AbortWithStatusJSON(httpCode, gin.H{
		"error":             code,
		"error_description": desc,
	})
}

// X generates a gin handler func based on the specified
// configuration and database.
func X(db authdb.Database, config Config) gin.HandlerFunc {
	scopeStr := ""
	for _, val := range config.Scope {
		scopeStr += val + " "
	}

	realmStr := ""
	for _, val := range config.Realms {
		realmStr += val + " "
	}

	return func(c *gin.Context) {
		c.Set("header-scopes", scopeStr)
		c.Set("header-realms", realmStr)
		tkStr := strings.Split(c.GetHeader("Authorization"), " ")
		if len(tkStr) != 2 {
			setError(c, ErrInvalid, "expected Authorization header")
			return
		}

		if tkStr[0] != "Bearer" {
			setError(c, ErrInvalid, "auth scheme must be Bearer")
			return
		}

		// Verify Access token exists.
		tkExists, err := db.Tokens().FindAccessByID(tkStr[1])
		if err != nil {
			setError(c, ErrToken, "access token expired/not found")
			return
		}

		// Ensure key is not expired.
		if (tkExists.CreatedAt + tkExists.TTL) < time.Now().Unix() {
			db.Tokens().DeleteAccessByID(tkStr[1])
			setError(c, ErrToken, "Access token expired / not found")
			return
		}

		// Verify associated user exists.
		userExists, err := db.Users().FindByID(tkExists.UserID)
		if err != nil {
			db.Tokens().DeleteAccessByID(tkStr[1])
			setError(c, ErrToken, "User not found")
			return
		}

		// Verify user has required realms.
		haveRealms := map[string]bool{}
		for _, realm := range userExists.Realms {
			haveRealms[realm] = true
		}

		for _, realm := range config.Realms {
			if !haveRealms[realm] {
				setError(c, ErrScope, "missing realms")
				return
			}
		}

		// Verify associated client exists.
		clientExists := authdb.ClientModel{
			ID: "0",
			Scope: []string{
				"dashboard", "users.update",
				"users.read", "users.delete",
				"clients.read", "clients.delete",
				"clients.create",
			},
		}

		if tkExists.ClientID != "0" {
			clientExists, err = db.Clients().FindByID(tkExists.ClientID)
			if err != nil {
				setError(c, ErrToken, "client not found")
				return
			}
		}

		// Verify token is authorized for this route.
		haveScope := map[string]bool{}
		for _, value := range clientExists.Scope {
			haveScope[value] = true
		}

		for _, value := range config.Scope {
			if !haveScope[value] {
				setError(c, ErrScope, "insufficient client scope")
				return
			}
		}

		// Write client, user, token to context.
		c.Set("user", userExists)
		c.Set("client", clientExists)
		c.Next()
	}
}
