package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetClientRoute reads a URL parameter specifying a client id and returns
// an error or data relevant to that client.
func GetClientRoute(db *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		client, err := db.ReadClient(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"error":  ErrClientNotFound,
			})
			return
		}

		// Found client.
		c.JSON(http.StatusOK, gin.H{
			"status":        "success",
			"id":            client.ID,
			"response_type": client.ResponseType,
			"redirect_uri":  client.RedirectURI,
			"scope":         client.Scope,
		})
	}
}
