package authapi

import (
	"github.com/gin-gonic/gin"
	"github.com/ufosc/OpenWebServices/pkg/authdb"
)

// APIController is an interface for retrieving gin middleware for
// each route in the authorization server.
type APIController interface {
	SignUpRoute() gin.HandlerFunc
	SignInRoute() gin.HandlerFunc
	VerifyEmailRoute() gin.HandlerFunc
	AuthClientRoute() gin.HandlerFunc

	AuthorizationRoute() gin.HandlerFunc
	TokenRoute() gin.HandlerFunc

	GetUserRoute() gin.HandlerFunc
	UpdateUserRoute() gin.HandlerFunc
	GetUsersRoute() gin.HandlerFunc
	ResetPwdRoute() gin.HandlerFunc

	GetClientRoute() gin.HandlerFunc
	CreateClientRoute() gin.HandlerFunc
	DeleteClientRoute() gin.HandlerFunc
	GetClientsRoute() gin.HandlerFunc

	DB() authdb.Database
	Stop() error
}

// DefaultAPIController implements APIController using authdb.
type DefaultAPIController struct {
	db     authdb.Database
	ms     MailSender
	secret string
}

// CreateAPIController creates an instance of APIController using uri and
// name as the MongoDB connection string and database name, respectively.
// addr is the email address to send verification emails from. secret is
// the random string for signing JWTs.
func CreateAPIController(uri, name, addr, secret string) (APIController, error) {
	cntrl := new(DefaultAPIController)
	cntrl.ms = NewMailSender(addr)
	db, err := authdb.NewDatabase(uri, name)
	if err != nil {
		return nil, err
	}
	if err := cntrl.ms.Start(1); err != nil {
		return nil, err
	}
	cntrl.db = db
	cntrl.secret = secret
	return cntrl, nil
}

// Stop the underlying database.
func (cntrl *DefaultAPIController) Stop() error {
	if err := cntrl.ms.Stop(); err != nil {
		return err
	}
	return cntrl.db.(*authdb.MongoDatabase).Stop()
}

// DB returns the underlying database interface.
func (cntrl *DefaultAPIController) DB() authdb.Database {
	return cntrl.db
}
