package main

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// GetClientRoute reads a URL parameter specifying a client id and returns
// an error or data relevant to that client.
func GetClientRoute(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		client, err := db.ReadClient(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Client not found",
			})
			return
		}

		// Found client.
		c.JSON(http.StatusOK, gin.H{
			"message":       "success",
			"id":            client.ID,
			"name":          client.Name,
			"description":   client.Description,
			"response_type": client.ResponseType,
			"redirect_uri":  client.RedirectURI,
			"scope":         client.Scope,
		})
	}
}

func CreateClientRoute(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// Extract user.
		userAny, ok := c.Get("user")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "User could not be found",
			})
			return
		}

		user, ok := userAny.(UserModel)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "User could not be found",
			})
			return
		}

		// Ensure description < 300 chars.
		if len(req.Description) > 300 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "description too long (max 300 chars)",
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

		// Validate RedirectURI.
		if !ValidateRedirectURI(req.RedirectURI) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid redirect_uri",
			})
			return
		}

		// Validate scope.
		if !ValidateScope(req.ResponseType, req.Scope) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid scope",
			})
			return
		}

		// Ensure user has client creation realm.
		hasRealm := false
		for _, value := range user.Realms {
			if value == "client.create" {
				hasRealm = true
				break
			}
		}

		if !hasRealm {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You are not authorized to create clients",
			})
			return
		}

		// Ensure name is unique.
		if _, err := db.ReadClientByName(req.Name); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Client name has already been registered",
			})
			return
		}

		// Create random private key.
		keyBytes := make([]byte, 256)
		if _, err := rand.Read(keyBytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error. Please try again later",
			})
			return
		}

		// Generate hash.
		hash, err := bcrypt.GenerateFromPassword(keyBytes,
			bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error. Please try again later",
			})
			return
		}

		// Create client & save to db.
		cid, err := db.CreateClient(ClientModel{
			Name:         req.Name,
			Description:  req.Description,
			ResponseType: req.ResponseType,
			RedirectURI:  req.RedirectURI,
			Scope:        req.Scope,
			Owner:        user.ID,
			Key:          string(hash),
			Created:      time.Now().Unix(),
			TTL:          5256000,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error. Please try again later",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "success",
			"client_id":     cid,
			"name":          req.Name,
			"description":   req.Description,
			"response_type": req.ResponseType,
			"redirect_uri":  req.RedirectURI,
			"scope":         req.Scope,
			"priv_key":      hex.EncodeToString(keyBytes),
		})
	}
}
