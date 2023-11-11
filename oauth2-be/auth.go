package main

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"
)

// AuthenticateUser returns a Gin middleware that handles authentication of a
// user. It checks whether the user JWT cookie is present, validates it, and
// handles any request errors.
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
		c.Next()
	}
}

func AuthenticateClient(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		tkStr := strings.Split(c.GetHeader("Authorization"), " ")
		if len(tkStr) != 2 {
			c.Header("WWW-Authenticate", "Basic realm=\"client\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrTokenUnauthorized,
			})
			return
		}

		if tkStr[0] != "Basic" {
			c.Header("WWW-Authenticate", "Basic realm=\"client\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrTokenUnauthorized,
			})
			return
		}

		// Extract client_id and secret key.
		decode, err := b64.StdEncoding.DecodeString(tkStr[1])
		if err != nil {
			c.Header("WWW-Authenticate", "Basic realm=\"client\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrTokenUnauthorized,
			})
			return
		}

		userPass := strings.Split(string(decode), ":")
		if len(userPass) != 2 {
			c.Header("WWW-Authenticate", "Basic realm=\"client\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrTokenUnauthorized,
			})
			return
		}

		// Verify.
		clientExists, err := db.ReadClient(userPass[0])
		if err != nil {
			c.Header("WWW-Authenticate", "Basic realm=\"client\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrTokenUnauthorized,
			})
			return
		}

		if !VerifyPassword(clientExists.Key, userPass[1]) {
			c.Header("WWW-Authenticate", "Basic realm=\"client\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrTokenUnauthorized,
			})
			return
		}

		// Check if client owner still exists.
		if _, err := db.ReadUser(clientExists.Owner); err != nil {
			db.DeleteClient(clientExists.ID)
			c.Header("WWW-Authenticate", "Basic realm=\"client\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrTokenUnauthorized,
			})
			return
		}

		// Set client context.
		c.Set("client", clientExists)
	}
}

// TODO: 'realm' in header response should reflect the values in 'scope'
func AuthenticateToken(db *Database, scope ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tkStr := strings.Split(c.GetHeader("Authorization"), " ")
		if len(tkStr) < 2 {
			c.Header("WWW-Authenticate", "Bearer realm=\"access_token\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrTokenUnauthorized,
			})
			return
		}

		if tkStr[0] != "Bearer" {
			c.Header("WWW-Authenticate", "Bearer realm=\"access_token\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrWrongTokenType,
			})
			return
		}

		// Verify key exists.
		tkExists, err := db.ReadAccessToken(tkStr[1])
		if err != nil {
			c.Header("WWW-Authenticate", "Bearer realm=\"access_token\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrTokenUnauthorized,
			})
			return
		}

		// Ensure key is not expired.
		if (tkExists.Created + tkExists.TTL) > time.Now().Unix() {
			c.Header("WWW-Authenticate", "Bearer realm=\"access_token\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrTokenExpired,
			})
			return
		}

		// Verify associated user exists.
		userExists, err := db.ReadUser(tkExists.UserID)
		if err != nil {
			c.Header("WWW-Authenticate", "Bearer realm=\"access_token\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrUserNotFound,
			})
			return
		}

		// Verify associated client exists.
		clientExists, err := db.ReadClient(tkExists.ClientID)
		if err != nil {
			c.Header("WWW-Authenticate", "Bearer realm=\"access_token\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrClientNotFound,
			})
			return
		}

		// Verify token is authorized for this route.
		expectedScope := map[string]bool{}
		for _, value := range scope {
			expectedScope[value] = true
		}

		for _, value := range clientExists.Scope {
			if expectedScope[value] {
				continue
			}
			c.Header("WWW-Authenticate", "Bearer realm=\"access_token\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrInsufficientPermission,
			})
			return
		}

		// Write client, user, token to context.
		c.Set("user", userExists)
		c.Set("client", clientExists)
		c.Set("token", tkExists)
		c.Next()
	}
}

// SignUpRoute returns the middleware for user sign up. It sends an email asking
// the user to verify their email address.
func SignupRoute(db *Database, ms MailSender) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// Validate Email.
		if !ValidateEmail(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrInvalidEmail,
			})
			return
		}

		// Ensure password is sufficiently strong.
		if err := ValidatePassword(req.Password); err != "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err,
			})
			return
		}

		// Ensure email is unique.
		if _, err := db.ReadUserByEmail(req.Email); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrEmailTaken,
			})
			return
		}

		// Ensure verification email hasn't already been sent.
		if _, err := db.ReadPendingUserByEmail(req.Email); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrEmailAlreadySent,
			})
			return
		}

		/* TODO: CAPTCHA
		   if !validateCaptcha(req.Captcha) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrInvalidCaptcha,
			})
			return
		}
		*/

		// Hash & salt password.
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password),
			bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrHashError,
			})
			return
		}

		// Create pending user instance.
		pendingUser := PendingUserModel{
			Email: req.Email,
			User: UserModel{
				"", req.Email, string(hash), req.FirstName,
				req.LastName, []string{}, 0, 0,
			},
			TTL: 1200,
		}

		// Save pending user to database.
		id, err := db.CreatePendingUser(pendingUser)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrDbFailure,
			})
			return
		}

		// Send verification email.
		if !ms.SendSignupVerification(id, pendingUser) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrEmailFailure,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "awaiting email verification",
		})
	}
}

func AuthorizeRoute(db *Database, config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.GIN_MODE == "debug" && config.FE_PROXY_PORT != "" {
			c.Redirect(http.StatusFound, "http://localhost:"+config.FE_PROXY_PORT)
			return
		}

		// TODO:
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

// SignInRoute returns the middleware for signing in. If successful, it
// assigns a JWT auth cookie to the user agent. If the user hasn't verified
// their address in over 3 months, it attempts to send a verification email.
func SigninRoute(db *Database, config Config, ms MailSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		// Extract JSON body.
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrMissingFields,
			})
			return
		}

		// Verify email exists.
		userExists, err := db.ReadUserByEmail(req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrIncorrectUserPass,
			})
			return
		}

		if !VerifyPassword(userExists.Password, req.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrIncorrectUserPass,
			})
			return
		}

		// Ensure email is still verified (i.e. within last 3 months).
		lastVerified := time.Unix(userExists.LastVerified, 0)
		if time.Since(lastVerified).Hours() > 2160 {

			// Check for pending sign-in verification.
			if _, err := db.ReadVerifyEmailSigninByEmail(req.Email); err == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "error",
					"error":  ErrEmailAlreadySent,
				})
				return
			}

			// Create pending sign-in verification instance.
			eid, err := db.CreateVerifyEmailSignin(VerifyEmailSigninModel{
				Email: req.Email,
				TTL:   1200,
			})

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "error",
					"error":  ErrDbFailure,
				})
				return
			}

			if !ms.SendSigninVerification(eid, req.Email) {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "error",
					"error":  ErrEmailFailure,
				})
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrVerifyEmail,
			})

			return
		}

		// Generate JWT.
		jwt := NewJWT(config.SECRET, userExists)
		if jwt == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrNewJWT,
			})
			return
		}

		// Assign JWT cookie.
		// TODO: Change localhost domain, consider SSL options.
		c.SetCookie("ows-jwt", jwt, 7200, "/", "localhost", false, true)
	}
}

// VerifyEmailRoute returns the middleware that consumes email verification
// requests. It deletes the underlying PendingUser or VerifyEmailSigninModel
// objects.
func VerifyEmailRoute(db *Database) gin.HandlerFunc {
	// TODO: Should serve user-friendly site instead of JSON.
	return func(c *gin.Context) {
		vtype := c.DefaultQuery("type", "")
		ref := c.DefaultQuery("ref", "")
		if (vtype != "signup" && vtype != "signin") || ref == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrInvalidURLParam,
			})
			return
		}

		// Handle sign up routine.
		if vtype == "signup" {
			pending, err := db.ReadPendingUser(ref)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "error",
					"error":  "bad or expired URL. Please try again",
				})
				return
			}

			// Delete pending user model.
			if err := db.DeletePendingUser(ref); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "error",
					"error":  ErrDbFailure,
				})
				return
			}

			// Update user creation dates.
			pending.User.Created = time.Now().Unix()
			pending.User.LastVerified = time.Now().Unix()

			// Sign up.
			if _, err := db.CreateUser(pending.User); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "error",
					"error":  ErrDbFailure,
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "success"})
			return
		}

		// Handle sign in.
		verif, err := db.ReadVerifyEmailSignin(ref)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "bad or expired URL. Please try again.",
			})
			return
		}

		// Delete pending email verification.
		if err := db.DeleteVerifyEmailSignin(ref); err != nil {
			fmt.Println("Error 1: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrDbFailure,
			})
			return
		}

		// Get underlying user.
		usr, err := db.ReadUserByEmail(verif.Email)
		if err != nil {
			fmt.Println("Error 2: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrDbFailure,
			})
			return
		}

		// Update LastVerified attribute.
		usr.LastVerified = time.Now().Unix()

		// Update user.
		if mod, err := db.UpdateUser(usr); err != nil || mod == 0 {
			fmt.Println("Error 3: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrDbFailure,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

func GrantRoute(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ResponseType string   `json:"response_type" binding:"required"`
			ClientID     string   `json:"client_id" binding:"required"`
			RedirectURI  string   `json:"redirect_uri" binding:"required"`
			Scope        []string `json:"scope" binding:"required"`
			State        string   `json:"state" binding:"required"`
		}

		// Get underlying user (passed in via AuthenticateUser route).
		userAny, ok := c.Get("user")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrUserNotFound,
			})
			return
		}

		// Cast to user model.
		user, ok := userAny.(UserModel)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrUserNotFound,
			})
			return
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
		if !ValidateScope(req.ResponseType, req.Scope) {
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
		client, err := db.ReadClient(req.ClientID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrClientNotFound,
			})
		}

		// Verify request redirect_uri matches client redirect_uri.
		if client.RedirectURI != req.RedirectURI {
			c.JSON(http.StatusUnauthorized, gin.H{
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
		sort.Strings(client.Scope)
		sort.Strings(req.Scope)
		if !reflect.DeepEqual(client.Scope, req.Scope) {
			redirect := fmt.Sprintf("%s?error=invalid_scope&state=%s",
				client.RedirectURI, req.State)
			c.Redirect(http.StatusFound, redirect)
			return
		}

		// Create implicit token.
		if client.ResponseType == "token" {
			token := Token{
				ClientID: client.ID,
				UserID:   user.ID,
				Created:  time.Now().Unix(),
				TTL:      1200,
			}

			// Save to db.
			id, err := db.CreateAccessToken(token)
			if err != nil {
				redirect := fmt.Sprintf("%s?error=server_error&state=%s",
					client.RedirectURI, req.State)
				c.Redirect(http.StatusFound, redirect)
				return
			}

			// Redirect user.
			uri := fmt.Sprintf("%s?access_token=%s&token_type=bearer&expires_in=1200&state=%s",
				client.RedirectURI, id, req.State)

			c.Redirect(http.StatusFound, uri)
			return
		}

		// Create authorization code.
		code := Token{
			ClientID: client.ID,
			UserID:   user.ID,
			Created:  time.Now().Unix(),
			TTL:      600,
		}

		// Save to DB.
		id, err := db.CreateAuthCode(code)
		if err != nil {
			redirect := fmt.Sprintf("%s?error=server_error&state=%s",
				client.RedirectURI, req.State)
			c.Redirect(http.StatusFound, redirect)
			return
		}

		// Redirect user.
		uri := fmt.Sprintf("%s?code=%s&state=%s", client.RedirectURI,
			id, req.State)

		c.Redirect(http.StatusFound, uri)
	}
}

func TokenRoute(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get underlying client (passed in via AuthenticateClient route).
		clientAny, ok := c.Get("client")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrClientNotFound,
			})
			return
		}

		// Cast to client model.
		client, ok := clientAny.(ClientModel)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrClientNotFound,
			})
			return
		}

		// Gather query parameters.
		grantType := c.DefaultQuery("grant_type", "")
		code := c.DefaultQuery("code", "")
		redirectUri := c.DefaultQuery("redirect_uri", "")
		clientID := c.DefaultQuery("client_id", "")
		refreshToken := c.DefaultQuery("refresh_token", "")

		// Validate grant type.
		if grantType != "refresh_token" && grantType != "authorization_code" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrInvalidGrantType,
			})
			return
		}

		if grantType == "refresh_token" {
			// Validate refresh token parameter.
			if refreshToken == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "error",
					"error":  ErrInvalidRefreshToken,
				})
				return
			}

			// Verify refresh token exists.
			token, err := db.ReadRefreshToken(refreshToken)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "error",
					"error":  ErrInvalidRefreshToken,
				})
				return
			}

			// Ensure token is not expired.
			if (token.Created + token.TTL) > time.Now().Unix() {
				db.DeleteRefreshToken(token.ID)
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "error",
					"error":  ErrTokenExpired,
				})
				return
			}

			// Ensure that the refresh token was issued to the
			// authenticated client.
			if token.ClientID != client.ID {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status": "error",
					"error":  ErrRefreshWrongClient,
				})
				return
			}

			// Ensure associated user still exists.
			if _, err := db.ReadUser(token.UserID); err != nil {
				db.DeleteRefreshToken(token.ID)
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "error",
					"error":  ErrUserNotFound,
				})
				return
			}

			// Create new access token.
			atoken := Token{
				ClientID: client.ID,
				UserID:   token.UserID,
				Created:  time.Now().Unix(),
				TTL:      1200,
			}

			// Save new access token to db.
			atokenID, err := db.CreateAccessToken(atoken)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "error",
					"error":  ErrDbFailure,
				})
				return
			}

			// Return token.
			c.JSON(http.StatusOK, gin.H{
				"status":        "success",
				"token":         atokenID,
				"token_type":    "bearer",
				"expires_in":    1200,
				"refresh_token": token.ID,
			})

			return
		}

		// Grant type : authorization_code.

		// Ensure required params are non-nil.
		if code == "" || redirectUri == "" || clientID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrMalformedURLParams,
			})
			return
		}

		// Validate code.
		codeExists, err := db.ReadAuthCode(code)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrTokenExpired,
			})
			return
		}

		// Ensure code has not expired.
		if (codeExists.Created + codeExists.TTL) > time.Now().Unix() {
			db.DeleteAuthCode(codeExists.ID)
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  ErrTokenExpired,
			})
			return
		}

		// Ensure client IDs match.
		if clientID != codeExists.ClientID {
			db.DeleteAuthCode(codeExists.ID)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrGrantWrongClient,
			})
			return
		}

		// Ensure client id exists.
		clientExists, err := db.ReadClient(clientID)
		if err != nil {
			db.DeleteAuthCode(codeExists.ID)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrClientNotFound,
			})
			return
		}

		// Ensure redirectURI is the same as client.
		if clientExists.RedirectURI != redirectUri {
			db.DeleteAuthCode(codeExists.ID)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrDifferentRedirectURI,
			})
			return
		}

		// Ensure userID still exists.
		if _, err := db.ReadUser(codeExists.UserID); err != nil {
			db.DeleteAuthCode(codeExists.ID)
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"error":  ErrUserNotFound,
			})
			return
		}

		// Create access token.
		atoken := Token{
			ClientID: client.ID,
			UserID:   codeExists.UserID,
			Created:  time.Now().Unix(),
			TTL:      1200,
		}

		aid, err := db.CreateAccessToken(atoken)
		if err != nil {
			db.DeleteAuthCode(codeExists.ID)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrDbFailure,
			})
			return
		}

		// Create refresh token.
		rtoken := Token{
			ClientID: client.ID,
			UserID:   codeExists.UserID,
			Created:  time.Now().Unix(),
			TTL:      5256000,
		}

		rid, err := db.CreateRefreshToken(rtoken)
		if err != nil {
			db.DeleteAccessToken(aid)
			db.DeleteAuthCode(codeExists.ID)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  ErrDbFailure,
			})
			return
		}

		// Return token.
		c.JSON(http.StatusOK, gin.H{
			"status":        "success",
			"access_token":  aid,
			"token_type":    "bearer",
			"expires_in":    1200,
			"refresh_token": rid,
		})
	}
}

func DeleteTokenRoute(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
