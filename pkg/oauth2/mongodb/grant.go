package mongodb

import (
	"context"
	"fmt"
	"github.com/ufosc/OpenWebServices/pkg/oauth2/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// CreateGrant saves an authorization grant to the database.
func (m *MongoDB) CreateGrant(g db.Grant) (string, error) {
	if m.isStopped {
		return "", fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("grants")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, g)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to create Grant")
	}

	return oid.String(), nil
}

// ReadGrant finds a grant in the database with the given ID.
func (m *MongoDB) ReadGrant(id string) (db.Grant, error) {
	if m.isStopped {
		return db.Grant{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("grants")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	var result db.Grant
	err := coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return db.Grant{}, err
	}

	return result, nil
}

// UpdateGrant updates the grant with the given ID, using the fields in grant g.
func (m *MongoDB) UpdateGrant(id string, g db.Grant) (db.Grant, error) {
	if m.isStopped {
		return db.Grant{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("grants")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	result, err := coll.UpdateOne(ctx, bson.D{{"_id", id}},
		bson.D{{"$set", g}})

	if err != nil {
		return db.Grant{}, err
	}

	if result.ModifiedCount != 1 {
		return db.Grant{}, fmt.Errorf("no grants modified")
	}

	return g, nil
}

// DeleteGrant deletes the grant with the given ID.
func (m *MongoDB) DeleteGrant(id string) (db.Grant, error) {
	if m.isStopped {
		return db.Grant{}, fmt.Errorf("mongoDB driver already stopped")
	}

	coll := m.client.Database(m.dbname).Collection("grants")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Find grant to return before deleting.
	grant, err := m.ReadGrant(id)
	if err != nil {
		return db.Grant{}, err
	}

	result, err := coll.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return db.Grant{}, err
	}

	if result.DeletedCount != 1 {
		return db.Grant{}, fmt.Errorf("no grants deleted")
	}

	return grant, nil
}
