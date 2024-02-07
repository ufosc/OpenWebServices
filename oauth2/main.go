package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ufosc/OpenWebServices/pkg/authapi"
	"github.com/ufosc/OpenWebServices/pkg/authmw"
	"net/http"
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
		AllowMethods:     []string{"POST, PUT, GET, DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API controller.
	api, err := authapi.CreateAPIController(config.MONGO_URI,
		config.DB_NAME, config.NOTIF_EMAIL_ADDR,
		config.WEBSMTP)

	if err != nil {
		panic(err)
	}

	defer api.Stop()

	// Auth.
	r.POST("/auth/signup", api.SignUpRoute())
	r.POST("/auth/signin", api.SignInRoute())
	r.GET("/auth/verify/:ref", api.VerifyEmailRoute())
	r.GET("/auth/token", api.TokenRoute())
	r.GET("/auth/authorize", authmw.A(api.DB()),
		api.AuthorizationRoute())

	// Resources.
	xEmpty := authmw.X(api.DB(), authmw.Config{})
	r.GET("/client/:id", api.GetClientRoute())
	r.GET("/user", xEmpty, api.GetUserRoute())

	r.PUT("/user", authmw.X(api.DB(), authmw.Config{
		Scope: []string{"users.modify"},
	}), api.UpdateUserRoute())

	r.GET("/users", authmw.X(api.DB(), authmw.Config{
		Scope:  []string{"users.read"},
		Realms: []string{"users.read"},
	}), api.GetUsersRoute())

	r.POST("/client", authmw.X(api.DB(), authmw.Config{
		Scope:  []string{"clients.create"},
		Realms: []string{"clients.create"},
	}), api.CreateClientRoute())

	r.GET("/clients", authmw.X(api.DB(), authmw.Config{
		Scope:  []string{"clients.read"},
		Realms: []string{"clients.read"},
	}), api.GetClientsRoute())

	r.DELETE("/client/:id", authmw.X(api.DB(), authmw.Config{
		Scope:  []string{"clients.delete"},
		Realms: []string{"clients.delete"},
	}), api.DeleteClientRoute())

	// Status.
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	r.Run("0.0.0.0:" + config.PORT)
}
