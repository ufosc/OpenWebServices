package mongodb

import (
	"context"
	"fmt"
	"github.com/ufosc/OpenWebServices/pkg/oauth2/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// CreateClient saves a client to the database.
func (m *MongoDB) CreateClient(c db.Client) (string, error) {
	if m.isStopped {
		return "", fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("clients")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, c)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to create Client")
	}

	return oid.String(), nil
}

// ReadClient finds a client in the database with the given ID.
func (m *MongoDB) ReadClient(id string) (db.Client, error) {
	if m.isStopped {
		return db.Client{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("clients")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	var result db.Client
	err := coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return db.Client{}, err
	}

	return result, nil
}

// UpdateClient updates the client with the given ID, using the fields in client c.
func (m *MongoDB) UpdateClient(id string, c db.Client) (db.Client, error) {
	if m.isStopped {
		return db.Client{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("clients")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	result, err := coll.UpdateOne(ctx, bson.D{{"_id", id}},
		bson.D{{"$set", c}})

	if err != nil {
		return db.Client{}, err
	}

	if result.ModifiedCount != 1 {
		return db.Client{}, fmt.Errorf("no clients modified")
	}

	return c, nil
}

// DeleteClient deletes the client with the given ID.
func (m *MongoDB) DeleteClient(id string) (db.Client, error) {
	if m.isStopped {
		return db.Client{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("clients")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Find client to return before deleting.
	client, err := m.ReadClient(id)
	if err != nil {
		return db.Client{}, err
	}

	result, err := coll.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return db.Client{}, err
	}

	if result.DeletedCount != 1 {
		return db.Client{}, fmt.Errorf("no clients deleted")
	}

	return client, nil
}
