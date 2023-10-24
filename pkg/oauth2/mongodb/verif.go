package mongodb

import (
	"context"
	"fmt"
	"github.com/ufosc/OpenWebServices/pkg/oauth2/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// CreateVerif saves a verification request to the database.
func (m *MongoDB) CreateVerif(v db.Verif) (string, error) {
	if m.isStopped {
		return "", fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("verifs")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, v)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to create Verif")
	}

	return oid.String(), nil
}

// ReadVerif finds a verification in the database with the given ID.
func (m *MongoDB) ReadVerif(id string) (db.Verif, error) {
	if m.isStopped {
		return db.Verif{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("verifs")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	var result db.Verif
	err := coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return db.Verif{}, err
	}

	return result, nil
}

// UpdateVerif updates the verif with the given ID, using the fields in verif v.
func (m *MongoDB) UpdateVerif(id string, v db.Verif) (db.Verif, error) {
	if m.isStopped {
		return db.Verif{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("verifs")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	result, err := coll.UpdateOne(ctx, bson.D{{"_id", id}},
		bson.D{{"$set", v}})

	if err != nil {
		return db.Verif{}, err
	}

	if result.ModifiedCount != 1 {
		return db.Verif{}, fmt.Errorf("no verifs modified")
	}

	return v, nil
}

// DeleteVerif deletes the verif with the given ID.
func (m *MongoDB) DeleteVerif(id string) (db.Verif, error) {
	if m.isStopped {
		return db.Verif{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("verifs")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Find verif to return before deleting.
	verif, err := m.ReadVerif(id)
	if err != nil {
		return db.Verif{}, err
	}

	result, err := coll.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return db.Verif{}, err
	}

	if result.DeletedCount != 1 {
		return db.Verif{}, fmt.Errorf("no verifs deleted")
	}

	return verif, nil
}
