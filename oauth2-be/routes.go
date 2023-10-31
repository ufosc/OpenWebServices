package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

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

		// Ensure email is unique.
		_, err := db.ReadUser(req.Email)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  ErrEmailTaken,
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
				"error": ErrHashError,
			})
			return
		}

		// Create pending user instance.
		pendingUser := PendingUserModel{
			User: UserModel{
				"", req.Email, string(hash), req.FirstName,
				req.LastName, 0, 0,
			},
			TTL: 20 * time.Minute,
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

func SigninRoute(db *Database, ms MailSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
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

		// Verify username exists.
		userExists, err := db.ReadUser(req.Username)
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
	}
}
