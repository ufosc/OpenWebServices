package authapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ufosc/OpenWebServices/pkg/authdb"
	"github.com/ufosc/OpenWebServices/pkg/common"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// SignUpRoute creates a new pending user and sends an email
// verification request.
func (cntrl *DefaultAPIController) SignUpRoute() gin.HandlerFunc {
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
				"error": "Missing required fields",
			})
			return
		}

		if len(req.FirstName) < 2 || len(req.LastName) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "first and/or last name are too short",
			})
			return
		}

		if len(req.FirstName) > 20 || len(req.LastName) > 20 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "first and/or last name are too long (> 20 chars)",
			})
			return
		}

		// Validate Email.
		if !common.ValidateEmail(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid email address",
			})
			return
		}

		// Ensure password is sufficiently strong.
		if err := common.ValidatePassword(req.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprint(err),
			})
			return
		}

		// Ensure email is unique.
		if _, err := cntrl.db.Users().FindByEmail(req.Email); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "An account already exists with this email",
			})
			return
		}

		// Ensure verification email hasn't already been sent.
		if _, err := cntrl.db.Users().FindPendingByEmail(req.Email); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Please verify your email address. If you have not received an email, " +
					"please check your spam folder and wait up to 10 minutes before trying again",
			})
			return
		}

		// Hash & salt password.
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password),
			bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error. Please try again later",
			})
			return
		}

		// Create pending user instance.
		pendingUser := authdb.PendingUserModel{
			Email: req.Email,
			User: authdb.UserModel{
				ID:        "",
				Email:     req.Email,
				Password:  string(hash),
				FirstName: req.FirstName,
				LastName:  req.LastName,
				Realms:    []string{},
				CreatedAt: 0,
			},
			CreatedAt: time.Now().Unix(),
			TTL:       600, // 10 minutes.
		}

		// Save pending user to database.
		id, err := cntrl.db.Users().CreatePending(pendingUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error. Please try again later",
			})
			return
		}

		// Send verification email.
		if !cntrl.SendVerification(id, pendingUser.Email) {
			cntrl.db.Users().DeletePendingByID(id)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error. Please try again later",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "awaiting email verification",
		})
	}
}

// SignInRoute authenticates a user and issues a JWT.
func (cntrl *DefaultAPIController) SignInRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		// Extract JSON body.
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing required fields",
			})
			return
		}

		// Verify email exists.
		userExists, err := cntrl.db.Users().FindByEmail(req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Incorrect username or password",
			})
			return
		}

		if !common.VerifyPassword(userExists.Password, req.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Incorrect username or password",
			})
			return
		}

		// Generate access token.
		tk, err := cntrl.db.Tokens().CreateAccess(authdb.TokenModel{
			ClientID:  "0",
			UserID:    userExists.ID,
			CreatedAt: time.Now().Unix(),
			TTL:       1200,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"token":   tk,
		})
	}
}

// VerifyEmailRoute consumes an email verification reference.
func (cntrl *DefaultAPIController) VerifyEmailRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		ref := c.Param("ref")
		if ref == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid URL",
			})
			return
		}

		pending, err := cntrl.db.Users().FindPendingByID(ref)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid URL",
			})
			return
		}

		// Delete pending user model.
		if err := cntrl.db.Users().DeletePendingByID(ref); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error. Please try again later",
			})
			return
		}

		// Update user creation dates.
		pending.User.CreatedAt = time.Now().Unix()

		// Sign up.
		if _, err := cntrl.db.Users().Create(pending.User); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error. Please try again later",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}
