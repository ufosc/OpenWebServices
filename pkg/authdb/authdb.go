package authdb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"sync/atomic"
)

// Database is an interface for accessing schema controllers.
type Database interface {
	Users() UserController
	Tokens() TokenController
	Clients() ClientController
}

// MongoState synchronizes database state and shares the MongoClient
// across the different schema controllers.
type MongoState struct {
	Name    string
	Client  *mongo.Client
	Stopped atomic.Bool
	Wg      sync.WaitGroup
}

// CollectionController is a controller struct template.
type CollectionController struct {
	state *MongoState
	coll  *mongo.Collection
}

// MongoDatabase implements database using a MongoDB connnection.
type MongoDatabase struct {
	state   MongoState
	clients ClientController
	tokens  TokenController
	users   UserController
}

// NewDatabase implements the Database interface using an underlying MongoDB
// driver connection, where uri is the connection URI and name is the db name.
func NewDatabase(uri, name string) (*MongoDatabase, error) {
	client, err := mongo.Connect(context.TODO(),
		options.Client().ApplyURI(uri))

	if err != nil {
		return nil, err
	}

	// Create DB controller.
	db := new(MongoDatabase)
	db.state.Name = name
	db.state.Client = client
	db.state.Stopped.Store(false)

	// Initialize controllers.
	clients, err := NewClientController(&db.state)
	if err != nil {
		return nil, err
	}
	db.clients = clients

	tokens, err := NewTokenController(&db.state)
	if err != nil {
		return nil, err
	}
	db.tokens = tokens

	users, err := NewUserController(&db.state)
	if err != nil {
		return nil, err
	}
	db.users = users

	return db, nil
}

// Stop the database.
func (db *MongoDatabase) Stop() error {
	// TODO: might be a good idea to add a mutex. What if Stop()
	// is called concurrently?
	if db.state.Stopped.Load() {
		return fmt.Errorf("db already stopped")
	}

	// Set the stop signal and wait for any workers to finish.
	db.state.Stopped.Store(true)
	db.state.Wg.Wait()

	// Disconnect the driver.
	if err := db.state.Client.Disconnect(context.TODO()); err != nil {
		return err
	}

	return nil
}

// Users returns the database user controller. Returns nil if closed.
func (db *MongoDatabase) Users() UserController {
	if db.state.Stopped.Load() {
		return nil
	}
	return db.users
}

// Tokens returns the database token controller. Returns nil if closed.
func (db *MongoDatabase) Tokens() TokenController {
	if db.state.Stopped.Load() {
		return nil
	}
	return db.tokens
}

// Clients returns the database client controller. Returns nil if closed.
func (db *MongoDatabase) Clients() ClientController {
	if db.state.Stopped.Load() {
		return nil
	}
	return db.clients
}
