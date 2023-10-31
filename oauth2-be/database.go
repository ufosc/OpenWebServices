package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database encapsulates a mongodb driver state.
type Database struct {
	dbname  string
	client  *mongo.Client
	stopped bool
}

// NewDatabase creates a new MongoDB driver instance that connects to the
// specified uri string and uses the dbname database.
func NewDatabase(uri, dbname string) (*Database, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	db := new(Database)
	db.dbname = dbname
	db.client = client
	db.stopped = false

	return db, nil
}

// Stop attempts to shut down the MongoDB driver connection.
func (db *Database) Stop() error {
	if db.stopped {
		return fmt.Errorf("db already stopped")
	}
	if err := db.client.Disconnect(context.TODO()); err != nil {
		return err
	}
	return nil
}

// ReadUser attempts to return a user by their email address.
func (db *Database) ReadUser(email string) (UserModel, error) {
	if db.stopped {
		return UserModel{}, fmt.Errorf("db has been stopped")
	}
	var user UserModel
	coll := db.client.Database(db.dbname).Collection("users")
	err := coll.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&user)
	if err != nil {
		return UserModel{}, err
	}

	return user, nil
}

func (db *Database) CreatePendingUser(usr PendingUserModel) (string, error) {
	if db.stopped {
		return "", fmt.Errorf("db has stopped")
	}

	coll := db.client.Database(db.dbname).Collection("pending_users")
	res, err := coll.InsertOne(context.TODO(), usr)
	if err != nil {
		return "", err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("internal server error")
	}

	return oid.String(), nil
}
