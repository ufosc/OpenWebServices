package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.POST("/auth/signin", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"username" binding:"required"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request",
			})
			return
		}

		if !validateUsername(req.Username) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "username must be a ufl.edu email address",
			})
			return
		}
	})

	r.POST("/auth/signup", func(c *gin.Context) {

	})

	r.POST("/auth/client", func(c *gin.Context) {

	})

	// When a user signs up, an email is sent asking them to click on a
	// button to verify their account. The button opens this route.
	r.GET("/auth/verify/:ref", func(c *gin.Context) {

	})

	r.GET("/auth/grant", func(c *gin.Context) {

	})

	r.GET("/auth/token", func(c *gin.Context) {

	})

	r.Run()
}

func validateUsername() {

}
