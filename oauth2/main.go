package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {

	// Web server.
	config := GetConfig()
	gin.SetMode(config.GIN_MODE)
	r := gin.Default()

	// Set up CORS.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Mail sender.
	ms := NewMailSender(config)
	ms.Start(1)
	defer ms.Stop()

	// Database.
	db, err := NewDatabase(config.MONGO_URI, config.DB_NAME)
	if err != nil {
		panic(err)
	}
	defer db.Stop()

	// Auth routes.
	r.POST("/auth/signup", SignupRoute(db, ms))
	r.POST("/auth/signin", SigninRoute(db, config, ms))
	r.POST("/auth/signout", SignoutRoute())
	r.GET("/auth/verify", VerifyEmailRoute(db))
	r.GET("/auth/grant", AuthenticateUser(db, config), GrantRoute(db))
	r.GET("/auth/token", AuthenticateClient(db), TokenRoute(db))
	r.DELETE("/auth/token/:id", AuthenticateClient(db), DeleteTokenRoute(db))

	// User API.
	r.GET("/user/:id", AuthenticateToken(db), func(c *gin.Context) {}) // Serve according to scope.
	r.PUT("/user/:id", AuthenticateToken(db, "user.modify"), func(c *gin.Context) {})

	// Client API: Clients do not get modified/updated. A new one must be
	// created.
	r.GET("/client/:id", GetClientRoute(db))
	r.POST("/client", AuthenticateToken(db, "client.create"), CreateClientRoute(db))

	r.Run()
}
