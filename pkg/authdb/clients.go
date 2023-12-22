package authdb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ClientModel is the client schema.
type ClientModel struct {
	ID           string   `bson:"_id,omitempty"`
	Name         string   `bson:"name"`
	Description  string   `bson:"description"`
	ResponseType string   `bson:"response_type"`
	RedirectURI  string   `bson:"redirect_uri"`
	Scope        []string `bson:"scope"`
	Owner        string   `bson:"owner"`
	Key          string   `bson:"key"`
	CreatedAt    int64    `bson:"createdAt"`
	TTL          int64    `bson:"expireAfterSeconds"`
}

// ClientController defines database operations for the OAuth2 client model.
type ClientController interface {
	FindByID(string) (ClientModel, error)
	FindByName(string) (ClientModel, error)
	Create(ClientModel) (string, error)
	DeleteByID(string) error
	Batch(n, skip int64) ([]ClientModel, error)
	Count() (int64, error)
}

// MongoClientController implements ClientController using MongoDB.
type MongoClientController CollectionController

// NewClientController creates a MongoDB client controller using the provided
// database state.
func NewClientController(state *MongoState) (ClientController, error) {
	if state == nil {
		return nil, ErrNilState
	}

	if state.Stopped.Load() {
		return nil, ErrClosed
	}

	ctrl := new(MongoClientController)
	ctrl.coll = state.Client.Database(state.Name).Collection("clients")
	ctrl.state = state

	return ctrl, nil
}

// FindById finds a client program by its given ID.
func (cc *MongoClientController) FindByID(id string) (ClientModel, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.coll == nil {
		return ClientModel{}, ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Extract primitive object ID.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ClientModel{}, err
	}

	// Find model.
	var client ClientModel
	err = cc.coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objID}}).Decode(&client)
	if err != nil {
		return ClientModel{}, err
	}

	return client, nil
}

// FindByName finds a client program by its advertised name.
func (cc *MongoClientController) FindByName(name string) (ClientModel, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.coll == nil {
		return ClientModel{}, ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Find model.
	var client ClientModel
	err := cc.coll.FindOne(context.TODO(), bson.D{{Key: "name", Value: name}}).Decode(&client)
	if err != nil {
		return ClientModel{}, err
	}

	return client, nil
}

// Create a client and save it to the database.
func (cc *MongoClientController) Create(client ClientModel) (string, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.coll == nil {
		return "", ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Insert.
	res, err := cc.coll.InsertOne(context.TODO(), client)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// DeleteByID deletes the client with the given id.
func (cc *MongoClientController) DeleteByID(id string) error {
	if cc.state == nil || cc.state.Stopped.Load() || cc.coll == nil {
		return ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Extract primitive object ID.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = cc.coll.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: objID}})
	return err
}

func (cc *MongoClientController) Batch(n, skip int64) ([]ClientModel, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.coll == nil {
		return []ClientModel{}, ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()
	cursor, err := cc.coll.Find(context.TODO(), bson.D{},
		options.Find().SetLimit(n).SetSkip(skip))

	if err != nil {
		return []ClientModel{}, err
	}

	result := []ClientModel{}
	if err := cursor.All(context.TODO(), &result); err != nil {
		return []ClientModel{}, err
	}

	return result, nil
}

func (cc *MongoClientController) Count() (int64, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.coll == nil {
		return -1, ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	count, err := cc.coll.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return -1, err
	}

	return count, nil
}
