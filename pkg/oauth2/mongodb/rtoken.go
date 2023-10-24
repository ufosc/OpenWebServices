package mongodb

import (
	"context"
	"fmt"
	"github.com/ufosc/OpenWebServices/pkg/oauth2/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// CreateRToken saves a refresh_token to the database.
func (m *MongoDB) CreateRToken(tk db.RToken) (string, error) {
	if m.isStopped {
		return "", fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("refresh_tokens")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, tk)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to create RToken")
	}

	return oid.String(), nil
}

// ReadRToken finds a refresh token in the database with the given ID.
func (m *MongoDB) ReadRToken(id string) (db.RToken, error) {
	if m.isStopped {
		return db.RToken{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("refresh_tokens")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	var result db.RToken
	err := coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return db.RToken{}, err
	}

	return result, nil
}

// UpdateRToken updates the refresh_token with the given ID, using the fields
// in RToken tk.
func (m *MongoDB) UpdateRToken(id string, tk db.RToken) (db.RToken, error) {
	if m.isStopped {
		return db.RToken{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("refresh_tokens")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	result, err := coll.UpdateOne(ctx, bson.D{{"_id", id}},
		bson.D{{"$set", tk}})

	if err != nil {
		return db.RToken{}, err
	}

	if result.ModifiedCount != 1 {
		return db.RToken{}, fmt.Errorf("no RTokens modified")
	}

	return tk, nil
}

// DeleteRToken deletes the refresh token with the given ID.
func (m *MongoDB) DeleteRToken(id string) (db.RToken, error) {
	if m.isStopped {
		return db.RToken{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("refresh_tokens")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Find RTokens to return before deleting.
	rtoken, err := m.ReadRToken(id)
	if err != nil {
		return db.RToken{}, err
	}

	result, err := coll.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return db.RToken{}, err
	}

	if result.DeletedCount != 1 {
		return db.RToken{}, fmt.Errorf("no RTokens deleted")
	}

	return rtoken, nil
}
