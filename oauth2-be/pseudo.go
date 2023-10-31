package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	config := GetConfig()
	gin.SetMode(config.GIN_MODE)
	r := gin.Default()

	r.POST("/auth/signup", func(c *gin.Context) {
		var req struct {
			Email     string `json:"email" binding:"required"`
			Password  string `json:"password" binding:"required"`
			FirstName string `json:"first_name" binding:"required"`
			LastName  string `json:"last_name" binding:"required"`
			Captcha   string `json:"captcha" binding:"required"`
		}

		// Extract JSON body.
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrMissingFields,
			})
			return
		}

		// Validate email.
		if !validateEmail(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrInvalidEmail,
			})
			return
		}

		// Ensure email is unique.
		if !isEmailUnique(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrEmailTaken,
			})
			return
		}

		// Ensure password is sufficiently strong.
		if err := validatePassword(req.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err,
			})
		}

		// Verify captcha token is correct.
		if !validateCaptcha(req.Captcha) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrInvalidCaptcha,
			})
			return
		}

		// Create pending user instance.
		pendingUser := PendingUserModel{
			User: UserModel{
				"", req.Email, req.Password, req.FirstName,
				req.LastName, 0, 0,
			},
			TTL: 20 * time.Minute,
		}

		// Save pending user to database.
		id, err := createPendingUser(pendingUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrDbFailure,
			})
			return
		}

		// Send verification email.
		if err := sendSignupVerification(id, pendingUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrEmailFailure,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	// TODO: Use access_token authentication for this route. Each client
	// should have a database attribute named 'scope'. If the user does
	// not have the realm for this scope (the scope should be privately
	// assigned), it should be rejected.
	r.POST("/auth/client", authenticateToken("create_client"), func(c *gin.Context) {
		var req struct {
			Name        string
			Type        string
			Description string
			RedirectURI string
			Scope       string
		}

		// Bind request body.
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrMissingFields,
			})
			return
		}

		// Validate type.
		if req.Type != "public" && req.Type != "confidential" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrInvalidClientType,
			})
			return
		}

		// Validate scope with respect to type.
		if !validateScope(req.Scope) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrInvalidScope,
			})
			return
		}

		// TODO: decide which scopes are available to what type of client.
		if req.Scope == "email" && req.Type == "public" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrInvalidPublicScope,
			})
			return
		}

		// Validate name and description (must be alphanumeric).
		if !validateClientInfo(req.Name, req.Description) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrInvalidClientInfo,
			})
			return
		}

		// Validate RedirectURI.
		if !validateRedirectURI(req.RedirectURI) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrInvalidRedirectURI,
			})
			return
		}

		// Name must be unique.
		if !isClientNameUnique(req.Name) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrClientNameTaken,
			})
			return
		}

		// RedirectURI must be unique.
		if !isRedirectURIUnique(req.RedirectURI) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrRedirectURITaken,
			})
			return
		}

		// Create client.
		newClient, err := createClient(ClientModel{
			Name:        req.Name,
			Type:        req.Type,
			Description: req.Description,
			RedirectURI: req.RedirectURI,
			Scope:       req.Scope,
			Created:     time.Now().Unix(),
			TTL:         7776000, // 90 days.
			Owner:       userID,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrDbFailure,
			})
			return
		}
	})

	r.GET("/auth/authorize", func(c *gin.Context) {
		// TODO: serve frontend.
		// TODO: do not give email address to public clients.
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "not implemented",
		})
	})

	r.POST("/auth/signin", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		// Verify username exists.
		userExists, err := readUser(req.Username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrIncorrectUserPass,
			})
			return
		}

		// Verify password matches registered username.
		if !verifyPassword(userExists, req.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrIncorrectUserPass,
			})
			return
		}

		// TODO: Ensure email is still verified.

		// Generate JWT.
		// TODO: check if the user is already logged in?
		jwt = NewJWT(userExists)

		// Assign JWT cookie.
		// TODO: Change localhost domain, consider SSL options.
		c.SetCookie("ows-jwt", jwt, 7200, "/", "localhost", false, true)
	})

	r.GET("/auth/grant", authenticateUser, func(c *gin.Context) {
		var req struct {
			ResponseType string   `json:"response_type" binding:"required"`
			ClientID     string   `json:"client_id" binding:"required"`
			RedirectURI  string   `json:"redirect_uri" binding:"required"`
			Scope        []string `json:"scope" binding:"required"`
			State        string   `json:"state" binding:"required"`
		}

		// TODO: The following error responses are client-facing.
		// Perhaps a dedicated error page should be served instead.

		// Extract JSON body.
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrMissingFields,
			})
			return
		}

		// Validate response type
		if req.ResponseType != "code" && req.ResponseType != "token" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrBadRespType,
			})
			return
		}

		// Validate scope.
		if !validateScope(req.Scope) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrInvalidScope,
			})
			return
		}

		// Validate state.
		if req.State == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrInvalidState,
			})
			return
		}

		// Verify client ID exists.
		client, err := readClient(req.ClientID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrDbFailure,
			})
			return
		}

		// Verify request redirect_uri matches client redirect_uri.
		if client.RedirectURI != req.RedirectURI {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrDifferentRedirectURI,
			})
			return
		}

		// Redirect to the client redirect_uri whenever possible (when
		// the integrity of redirect_uri can be verified).

		// Verify request response_type matches client configuration.
		if client.ResponseType != req.ResponseType {
			redirect := fmt.Sprintf("%s?error=invalid_request&state=%s",
				client.RedirectURI, req.State)
			c.Redirect(http.StatusFound, redirect)
			return
		}

		// Verify request scope matches client configuration.
		if client.Scope != req.Scope {
			redirect := fmt.Sprintf("%s?error=invalid_scope&state=%s",
				client.RedirectURI, req.State)
			c.Redirect(http.StatusFound, redirect)
			return
		}

		// Generate tokens.
		resparams := ""
		if client.ResponseType == "code" {
			token := NewAuthCode(client)
			resparams = fmt.Sprintf("?code=%s&state=%s",
				token, req.State)

			createGrantToken(token, client.ClientID, client.RedirectURI, username)
		}

		if client.ResponseType == "token" {
			token := NewImplicitToken(client)
			resparams = fmt.Sprintf("?access_token=%s&token_type=bearer&state=%s",
				token, req.State)

			createImplicitToken(token, client.ClientID, username)
		}

		// Redirect.
		c.Redirect(http.StatusFound, client.RedirectURI+resparams)
	})

	r.GET("/auth/token", authenticateClient, func(c *gin.Context) {
		grantType := c.DefaultQuery("grant_type", "")
		code := c.DefaultQuery("code", "")
		if (grantType != "authorization_code" && grantType !=
			"refresh_token") || code == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrMalformedURLParams,
			})
			return
		}

		if grantType == "authorization_code" {
			redirectURI := c.DefaultQuery("redirect_uri", "")
			urlClientID := c.DefaultQuery("client_id", "")
			if redirectURI == "" || urlClientID == "" {
				deleteGrantToken(code)
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "error",
					"error":  ErrMalformedURLParams,
				})
				return
			}

			tk, err := readGrantToken(code)
			if err != nil {
				deleteGrantToken(code)
				c.JSON(http.StatusUnauthorized, gin.H{
					"status": "error",
					"error":  ErrGrantTokenNotFound,
				})
				return
			}

			if tk.RedirectURI != redirectURI || tk.ClientID !=
				urlClientID || tk.ClientID != clientID {
				deleteGrantToken(code)
				c.JSON(http.StatusUnauthorized, gin.H{
					"status": "error",
					"error":  ErrGrantWrongClient,
				})
				return
			}

			deleteGrantToken(code)

			// TODO: The authorization server MUST maintain the
			// binding between a refresh token and the client to
			// whom it was issued.
			rtoken, err := createRefreshToken( /* todo */ )
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "error",
					"error":  ErrCreateRefreshToken,
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":        "success",
				"access_token":  rtoken.AccessToken,
				"token_type":    "bearer",
				"expires_in":    1200,
				"refresh_token": rtoken.RefreshToken,
			})

			return
		}

		// grantType == "refresh_token"

		tk, err := readRefreshToken(code)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrInvalidRefreshToken,
			})
			return
		}

		if tk.ClientID != clientID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status", "error",
				"error": ErrRefreshWrongClient,
			})
			return
		}

		deleteRefreshToken(tk)
		rtoken, err := createRefreshToken( /* todo */ )
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrCreateRefreshToken,
			})
			return
		}

		// Return refresh token.
		c.JSON(http.StatusOK, gin.H{
			"status":        "success",
			"access_token":  rtoken.AccessToken(),
			"token_type":    "bearer",
			"expires_in":    1200,
			"refresh_token": rtoken.RefreshToken(),
		})
	})

	r.GET("/auth/verify", func(c *gin.Context) {
		vtype := c.DefaultQuery("type", "")
		ref := c.DefaultQuery("ref", "")
		if (vtype != "signup" && vtype != "signin") || ref == "" {
			// TODO: serve error page.
			return
		}

		if vtype == "signup" {
			// TODO: when creating users, finalize the last_verified
			// and created attributes.
			return
		}

		// vtype == "signin"

		// TODO: verify email route.
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "not implemented",
		})
	})

	r.Run()
}
