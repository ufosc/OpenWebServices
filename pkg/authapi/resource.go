package authapi

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/ufosc/OpenWebServices/pkg/authdb"
	"github.com/ufosc/OpenWebServices/pkg/common"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// GetUserRoute returns user information based on the permissions
// defined by the client's scope.
func (cntrl *DefaultAPIController) GetUserRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "not implemented",
		})
	}
}

// UpdateUserRoute updates user information.
func (cntrl *DefaultAPIController) UpdateUserRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "not implemented",
		})
	}
}

// ForgotPassword sends the user an email to change their password.
func (cntrl *DefaultAPIController) ResetPwdRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "not implemented",
		})
	}
}

// GetClientRoute returns public information about a client.
func (cntrl *DefaultAPIController) GetClientRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.Param("id")
		clientExists, err := cntrl.db.Clients().FindByID(clientID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "client not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":            clientExists.ID,
			"name":          clientExists.Name,
			"description":   clientExists.Description,
			"response_type": clientExists.ResponseType,
			"redirect_uri":  clientExists.RedirectURI,
			"scope":         clientExists.Scope,
		})
	}
}

// CreateCreateClient returns the gin middleware for registering a new client.
// It expects user authentication (clients are owned by users).
func (cntrl *DefaultAPIController) CreateClientRoute() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Request body.
		var req struct {
			Name         string   `json:"name" binding:"required"`
			Description  string   `json:"description" binding:"required"`
			ResponseType string   `json:"response_type" binding:"required"`
			RedirectURI  string   `json:"redirect_uri" binding:"required"`
			Scope        []string `json:"scope" binding:"required"`
		}

		// Extract JSON body.
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing required fields",
			})
			return
		}

		// Get underlying user (from middleware).
		userAny, ok := c.Get("user")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "User not found",
			})
			return
		}

		// Cast to user model.
		user, ok := userAny.(authdb.UserModel)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "User not found",
			})
			return
		}

		// User must have client creation realm.
		hasRealm := false
		for _, v := range user.Realms {
			if v == "client.create" {
				hasRealm = true
				break
			}
		}

		if !hasRealm {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authorized to create clients",
			})
			return
		}

		// Validate response type.
		if req.ResponseType != "code" && req.ResponseType != "token" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "response_type must be 'code' or 'token'",
			})
			return
		}

		// Validate redirect uri.
		if !common.ValidateRedirectURI(req.RedirectURI) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid redirect_uri",
			})
			return
		}

		// Validate scope.
		if !common.ValidateScope(req.ResponseType, req.Scope) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid scope",
			})
			return
		}

		// Ensure name and description are not too long.
		if len(req.Name) > 12 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "name cannot be longer than 12 characters",
			})
			return
		}

		if len(req.Description) > 150 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "description cannot be longer than 150 characters",
			})
			return
		}

		// Ensure name doesn't already exist.
		if _, err := cntrl.db.Clients().FindByName(req.Name); err != mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "client name already registered",
			})
			return
		}

		// Generate random key.
		privateKey, _, _, err := elliptic.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error. Please try again later",
			})
			return
		}

		pkey := base64.StdEncoding.EncodeToString(privateKey)

		// Hash random key.
		pkeyHash, err := bcrypt.GenerateFromPassword([]byte(pkey),
			bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error. Please try again later",
			})
			return
		}

		// Create client model.
		client := authdb.ClientModel{
			ID:           "",
			Name:         req.Name,
			Description:  req.Description,
			ResponseType: req.ResponseType,
			RedirectURI:  req.RedirectURI,
			Scope:        req.Scope,
			Owner:        user.ID,
			Key:          string(pkeyHash),
			Created:      time.Now().Unix(),
			TTL:          7890000, // 3 months.
		}

		id, err := cntrl.db.Clients().Create(client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"id":      id,
			"pkey":    pkey,
		})
	}
}

// DeleteClientRoute returns the gin middleware for deleting a client.
func (cntrl *DefaultAPIController) DeleteClientRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "not implemented",
		})
	}
}