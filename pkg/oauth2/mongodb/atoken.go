package mongodb

import (
	"context"
	"fmt"
	"github.com/ufosc/OpenWebServices/pkg/oauth2/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// CreateAToken saves an access token to the database.
func (m *MongoDB) CreateAToken(tk db.AToken) (string, error) {
	if m.isStopped {
		return "", fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("access_tokens")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, tk)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to create AToken")
	}

	return oid.String(), nil
}

// ReadAToken finds an access token in the database with the given ID.
func (m *MongoDB) ReadAToken(id string) (db.AToken, error) {
	if m.isStopped {
		return db.AToken{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("access_tokens")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	var result db.AToken
	err := coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return db.AToken{}, err
	}

	return result, nil
}

// UpdateAToken updates the token with the given ID, using the fields in AToken tk.
func (m *MongoDB) UpdateAToken(id string, tk db.AToken) (db.AToken, error) {
	if m.isStopped {
		return db.AToken{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("access_tokens")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	result, err := coll.UpdateOne(ctx, bson.D{{"_id", id}},
		bson.D{{"$set", tk}})

	if err != nil {
		return db.AToken{}, err
	}

	if result.ModifiedCount != 1 {
		return db.AToken{}, fmt.Errorf("no ATokens modified")
	}

	return tk, nil
}

// DeleteAToken deletes the access token with the given ID.
func (m *MongoDB) DeleteAToken(id string) (db.AToken, error) {
	if m.isStopped {
		return db.AToken{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("access_tokens")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Find AToken to return before deleting.
	atoken, err := m.ReadAToken(id)
	if err != nil {
		return db.AToken{}, err
	}

	result, err := coll.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return db.AToken{}, err
	}

	if result.DeletedCount != 1 {
		return db.AToken{}, fmt.Errorf("no ATokens deleted")
	}

	return atoken, nil
}
