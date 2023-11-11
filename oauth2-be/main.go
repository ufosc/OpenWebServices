package main

import "github.com/gin-gonic/gin"

func main() {

	// Web server.
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

	// Auth routes.
	r.POST("/auth/signup", SignupRoute(db, ms))
	r.GET("/auth/authorize", AuthorizeRoute(db, config))
	r.POST("/auth/signin", SigninRoute(db, config, ms))
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
