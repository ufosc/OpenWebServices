package authdb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
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

	initIndices(db)
	return db, nil
}

// initIndices initializes database indices.
func initIndices(db *MongoDatabase) {
	index := func(ttl int32) mongo.IndexModel {
		return mongo.IndexModel{
			Keys:    bson.M{"createdAt": 1},
			Options: options.Index().SetExpireAfterSeconds(ttl),
		}
	}

	// Gather database collections.
	clicol := db.state.Client.Database(db.state.Name).Collection("clients")
	refcol := db.state.Client.Database(db.state.Name).Collection("refresh_tokens")
	acccol := db.state.Client.Database(db.state.Name).Collection("access_tokens")
	autcol := db.state.Client.Database(db.state.Name).Collection("auth_tokens")
	pencol := db.state.Client.Database(db.state.Name).Collection("pending_users")

	// Apply indices.
	_, err := clicol.Indexes().CreateOne(context.TODO(), index(7890000))
	if err != nil {
		fmt.Println("unable to apply TTL to client collection:", err)
		os.Exit(1)
	}

	_, err = refcol.Indexes().CreateOne(context.TODO(), index(5256000))
	if err != nil {
		fmt.Println("unable to apply TTL to refresh token collection:", err)
		os.Exit(1)
	}

	_, err = acccol.Indexes().CreateOne(context.TODO(), index(1200))
	if err != nil {
		fmt.Println("unable to apply TTL to access token collection:", err)
		os.Exit(1)
	}

	_, err = autcol.Indexes().CreateOne(context.TODO(), index(600))
	if err != nil {
		fmt.Println("unable to apply TTL to grant token collection:", err)
		os.Exit(1)
	}

	_, err = pencol.Indexes().CreateOne(context.TODO(), index(600))
	if err != nil {
		fmt.Println("unable to apply TTL to pending_users collection:", err)
		os.Exit(1)
	}

	// Create a custom identifier index for tokens and verification
	// emails. Default indices are not cryptographically random.
	_, err = refcol.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.M{"ID": 1},
	})

	if err != nil {
		fmt.Println("cannot apply index to refresh_token collection", err)
		os.Exit(1)
	}

	_, err = acccol.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.M{"ID": 1},
	})

	if err != nil {
		fmt.Println("cannot apply index to access_token collection", err)
		os.Exit(1)
	}

	_, err = autcol.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.M{"ID": 1},
	})

	if err != nil {
		fmt.Println("cannot apply index to auth_token collection", err)
		os.Exit(1)
	}

	_, err = pencol.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.M{"ID": 1},
	})

	if err != nil {
		fmt.Println("cannot apply index to pending_users collection", err)
		os.Exit(1)
	}
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
