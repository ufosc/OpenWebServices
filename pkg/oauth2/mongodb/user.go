package mongodb

import (
	"context"
	"fmt"
	"github.com/ufosc/OpenWebServices/pkg/oauth2/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// CreateUser saves a user to the database. It is assumed that the Created,
// LastLogin, and LastVerified attributes are already written in the user u.
func (m *MongoDB) CreateUser(u db.User) (string, error) {
	if m.isStopped {
		return "", fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("users")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, u)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to create User")
	}

	return oid.String(), nil
}

// ReadUser finds a user in the database with the given ID.
func (m *MongoDB) ReadUser(id string) (db.User, error) {
	if m.isStopped {
		return db.User{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("users")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	var result db.User
	err := coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return db.User{}, err
	}

	return result, nil
}

// UpdateUser updates the user with the given ID, using the fields in user u.
func (m *MongoDB) UpdateUser(id string, u db.User) (db.User, error) {
	if m.isStopped {
		return db.User{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("users")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	result, err := coll.UpdateOne(ctx, bson.D{{"_id", id}},
		bson.D{{"$set", u}})

	if err != nil {
		return db.User{}, err
	}

	if result.ModifiedCount != 1 {
		return db.User{}, fmt.Errorf("no users modified")
	}

	return u, nil
}

// DeleteUser deletes the user with the given ID.
func (m *MongoDB) DeleteUser(id string) (db.User, error) {
	if m.isStopped {
		return db.User{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("users")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Find user to return before deleting.
	user, err := m.ReadUser(id)
	if err != nil {
		return db.User{}, err
	}

	result, err := coll.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return db.User{}, err
	}

	if result.DeletedCount != 1 {
		return db.User{}, fmt.Errorf("no users deleted")
	}

	return user, nil
}
