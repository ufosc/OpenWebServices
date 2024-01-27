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
		config.WEBSMTP, config.SECRET)

	if err != nil {
		panic(err)
	}

	defer api.Stop()

	// Auth routes.
	r.POST("/auth/signup", api.SignUpRoute())
	r.POST("/auth/signin", api.SignInRoute())
	r.GET("/auth/verify/:ref", api.VerifyEmailRoute())
	r.POST("/auth/client", api.AuthClientRoute())
	r.GET("/auth/token", api.TokenRoute())

	r.GET("/auth/authorize", authmw.AuthenticateUser(config.SECRET,
		api.DB()), api.AuthorizationRoute())

	// Resource API.
	r.GET("/client/:id", api.GetClientRoute())

	r.GET("/user", authmw.AuthenticateBearer(config.SECRET, api.DB(),
		[]string{}, []string{}), api.GetUserRoute())

	r.PUT("/user", authmw.AuthenticateBearer(config.SECRET, api.DB(),
		[]string{}, []string{"modify"}), api.UpdateUserRoute())

	r.GET("/users", authmw.AuthenticateBearer(config.SECRET, api.DB(),
		[]string{"users.read"}, []string{}), api.GetUsersRoute())

	r.POST("/client", authmw.AuthenticateUser(config.SECRET, api.DB(),
		"clients.create"), api.CreateClientRoute())

	r.GET("/clients", authmw.AuthenticateBearer(config.SECRET, api.DB(),
		[]string{"clients.read"}, []string{}), api.GetClientsRoute())

	r.DELETE("/client/:id", authmw.AuthenticateUser(config.SECRET, api.DB()),
		api.DeleteClientRoute())

	// Status check.
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	r.Run("0.0.0.0:" + config.PORT)
}
