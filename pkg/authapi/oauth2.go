package authapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ufosc/OpenWebServices/pkg/authdb"
	"github.com/ufosc/OpenWebServices/pkg/authmw"
	"github.com/ufosc/OpenWebServices/pkg/common"
	"net/http"
	"net/url"
	"time"
)

// AuthorizationRoute returns the middleware for the Oauth2 authorize route.
func (cntrl *DefaultAPIController) AuthorizationRoute() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get underlying user.
		userAny, _ := c.Get("user")
		user, _ := userAny.(authdb.UserModel)

		// Gather query parameters.
		responseType := c.DefaultQuery("response_type", "")
		clientID := c.DefaultQuery("client_id", "")
		redirectURI := c.DefaultQuery("redirect_uri", "")
		state := c.DefaultQuery("state", "")

		// Validate response type
		if responseType != "code" && responseType != "token" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":             "invalid_request",
				"error_description": "response_type must be 'code' or 'token'",
			})
			return
		}

		// Validate state.
		if state == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":             "invalid_request",
				"error_description": "state parameter cannot be empty string",
			})
			return
		}

		// Verify client ID exists.
		client, err := cntrl.db.Clients().FindByID(clientID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":             "not_found",
				"error_description": "client ID not found",
			})
			return
		}

		redirectDecoded, err := url.QueryUnescape(redirectURI)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":             "invalid_request",
				"error_description": "redirect_uri is invalid",
			})
			return
		}

		// Verify request redirect_uri matches client redirect_uri.
		if client.RedirectURI != redirectDecoded {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":            "invalid_request",
				"error_descriptor": "wrong redirect_uri",
			})
			return
		}

		// Redirect to the client redirect_uri whenever possible (when
		// the integrity of redirect_uri can be verified).

		// Verify request response_type matches client configuration.
		if client.ResponseType != responseType {
			redirect := fmt.Sprintf("%s?error=invalid_request&state=%s",
				client.RedirectURI, state)
			c.Redirect(http.StatusFound, redirect)
			return
		}

		// Create implicit token.
		if client.ResponseType == "token" {
			token := authdb.TokenModel{
				ID:        common.UUID(),
				ClientID:  client.ID,
				UserID:    user.ID,
				CreatedAt: time.Now().Unix(),
				TTL:       1200,
			}

			// Save to db.
			id, err := cntrl.db.Tokens().CreateAccess(token)
			if err != nil {
				redirect := fmt.Sprintf("%s?error=server_error&state=%s",
					client.RedirectURI, state)
				c.Redirect(http.StatusFound, redirect)
				return
			}

			// Redirect user.
			uri := fmt.Sprintf("%s?access_token=%s&token_type=bearer&expires_in=1200&state=%s",
				client.RedirectURI, id, state)

			c.Redirect(http.StatusFound, uri)
			return
		}

		// Create authorization code.
		code := authdb.TokenModel{
			ID:        common.UUID(),
			ClientID:  client.ID,
			UserID:    user.ID,
			CreatedAt: time.Now().Unix(),
			TTL:       600,
		}

		// Save to DB.
		id, err := cntrl.db.Tokens().CreateAuth(code)
		if err != nil {
			redirect := fmt.Sprintf("%s?error=server_error&state=%s",
				client.RedirectURI, state)
			c.Redirect(http.StatusFound, redirect)
			return
		}

		// Redirect user.
		uri := fmt.Sprintf("%s?code=%s&state=%s", client.RedirectURI,
			id, state)

		c.Redirect(http.StatusFound, uri)
	}
}

// TokenRoute returns the gin middleware for the Oauth2 token route.
func (cntrl *DefaultAPIController) TokenRoute() gin.HandlerFunc {
	return func(c *gin.Context) {

		grantType := c.DefaultQuery("grant_type", "")
		if grantType == "authorization_code" {
			cntrl.handleAuthCode(c)
			return
		}

		if grantType == "refresh_token" {
			cntrl.handleRefreshToken(c)
			return
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid grant_type",
		})
	}
}

func (cntrl *DefaultAPIController) handleAuthCode(c *gin.Context) {
	authmw.B(cntrl.db)(c)
	if c.IsAborted() {
		return
	}

	// Get underlying client.
	clientAny, ok := c.Get("client")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "not_found",
			"error_description": "client ID not found",
		})
		return
	}

	// Cast to client model.
	client, ok := clientAny.(authdb.ClientModel)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "not_found",
			"error_description": "client ID not found",
		})
		return
	}

	grantType := c.DefaultQuery("grant_type", "")
	if grantType != "authorization_code" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": "expected authorization_code grant type",
		})
		return
	}

	code := c.DefaultQuery("code", "")
	redirectUri := c.DefaultQuery("redirect_uri", "")
	clientID := c.DefaultQuery("client_id", "")

	// Ensure required params are non-nil.
	if code == "" || redirectUri == "" || clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": "Malformed or missing URL parameters",
		})
		return
	}

	// Validate code.
	codeExists, err := cntrl.db.Tokens().FindAuthByID(code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":             "not_found",
			"error_description": "Token expired or could not be found",
		})
		return
	}

	// Ensure code has not expired.
	if (codeExists.CreatedAt + codeExists.TTL) < time.Now().Unix() {
		cntrl.db.Tokens().DeleteAuthByID(codeExists.ID)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":             "not_found",
			"error_description": "Token expired or could not be found",
		})
		return
	}

	// Ensure client IDs match.
	if clientID != codeExists.ClientID {
		cntrl.db.Tokens().DeleteAuthByID(codeExists.ID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "unauthorized",
			"error_description": "Token was not issued to this client",
		})
		return
	}

	// Ensure client id exists.
	clientExists, err := cntrl.db.Clients().FindByID(clientID)
	if err != nil {
		cntrl.db.Tokens().DeleteAuthByID(codeExists.ID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "not_found",
			"error_description": "The client associated with this token could not be found",
		})
		return
	}

	// Ensure redirectURI is the same as client.
	if clientExists.RedirectURI != redirectUri {
		cntrl.db.Tokens().DeleteAuthByID(codeExists.ID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": "redirect_uri and client-registered redirect_uri do not match",
		})
		return
	}

	// Ensure userID still exists.
	if _, err := cntrl.db.Users().FindByID(codeExists.UserID); err != nil {
		cntrl.db.Tokens().DeleteAuthByID(codeExists.ID)
		c.JSON(http.StatusNotFound, gin.H{
			"error":             "not_found",
			"error_description": "The user associated with this token could not be found",
		})
		return
	}

	// Create access token.
	atoken := authdb.TokenModel{
		ID:        common.UUID(),
		ClientID:  client.ID,
		UserID:    codeExists.UserID,
		CreatedAt: time.Now().Unix(),
		TTL:       1200,
	}

	aid, err := cntrl.db.Tokens().CreateAccess(atoken)
	if err != nil {
		cntrl.db.Tokens().DeleteAuthByID(codeExists.ID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "internal_server_error",
			"error_description": "Internal server error. Please try again later",
		})
		return
	}

	// Create refresh token.
	rtoken := authdb.TokenModel{
		ID:        common.UUID(),
		ClientID:  client.ID,
		UserID:    codeExists.UserID,
		CreatedAt: time.Now().Unix(),
		TTL:       5256000,
	}

	rid, err := cntrl.db.Tokens().CreateRefresh(rtoken)
	if err != nil {
		cntrl.db.Tokens().DeleteAccessByID(aid)
		cntrl.db.Tokens().DeleteAuthByID(codeExists.ID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "internal_server_error",
			"error_description": "Internal server error. Please try again later",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "success",
		"access_token":  aid,
		"token_type":    "bearer",
		"expires_in":    1200,
		"refresh_token": rid,
	})
}

func (cntrl *DefaultAPIController) handleRefreshToken(c *gin.Context) {
	authmw.B(cntrl.db)(c)
	if c.IsAborted() {
		return
	}

	// Get underlying client.
	clientAny, ok := c.Get("client")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "not_found",
			"error_description": "Client not found",
		})
		return
	}

	// Cast to client model.
	client, ok := clientAny.(authdb.ClientModel)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "not_found",
			"error_description": "Client not found",
		})
		return
	}

	grantType := c.DefaultQuery("grant_type", "")
	if grantType != "refresh_token" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": "expected refresh_token grant type",
		})
		return
	}

	refreshToken := c.DefaultQuery("refresh_token", "")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": "invalid refresh token",
		})
		return
	}

	// Verify refresh token exists.
	token, err := cntrl.db.Tokens().FindRefreshByID(refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "not_found",
			"error_description": "Refresh token expired or could not be found",
		})
		return
	}

	// Ensure token is not expired.
	if (token.CreatedAt + token.TTL) > time.Now().Unix() {
		cntrl.db.Tokens().DeleteRefreshByID(token.ID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "not_found",
			"error_description": "Refresh token expired or could not be found",
		})
		return
	}

	// Ensure that the refresh token was issued to the
	// authenticated client.
	if token.ClientID != client.ID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":             "not_found",
			"error_description": "Refresh token was not issued to this client",
		})
		return
	}

	// Ensure associated user still exists.
	if _, err := cntrl.db.Users().FindByID(token.UserID); err != nil {
		cntrl.db.Tokens().DeleteRefreshByID(token.ID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "not_found",
			"error_description": "the user associated with this token could not be found",
		})
		return
	}

	// Create new access token.
	atoken := authdb.TokenModel{
		ID:        common.UUID(),
		ClientID:  client.ID,
		UserID:    token.UserID,
		CreatedAt: time.Now().Unix(),
		TTL:       1200,
	}

	// Save new access token to db.
	atokenID, err := cntrl.db.Tokens().CreateAccess(atoken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "internal_server_error",
			"error_description": "internal server error. Please try again later",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "success",
		"token":         atokenID,
		"token_type":    "bearer",
		"expires_in":    1200,
		"refresh_token": token.ID,
	})
}
