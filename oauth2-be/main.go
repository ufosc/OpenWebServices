package main

import "github.com/gin-gonic/gin"

func main() {

	// Server.
	config := GetConfig()
	gin.SetMode(config.GIN_MODE)
	r := gin.Default()

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

	// Routes.
	r.POST("/auth/signup", SignupRoute(db, ms))
	r.GET("/auth/authorize", AuthorizeRoute(db, config))
	r.POST("/auth/signin", SigninRoute(db, config, ms))
	r.GET("/auth/verify", VerifyEmailRoute(db))
	r.GET("/auth/grant", AuthenticateUser(db, config), GrantRoute(db))
	r.GET("/auth/token", AuthenticateClient(db), TokenRoute(db))

	// User API.
	r.GET("/user/:id", AuthenticateToken(db), func(c *gin.Context) {}) // Serve according to scope.
	r.PUT("/user/:id", AuthenticateToken(db), func(c *gin.Context) {})

	// Client API.
	r.GET("/client/:id", func(c *gin.Context) {})
	r.POST("/client", AuthenticateToken(db), func(c *gin.Context) {})

	r.Run()
}
