package authdb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// The Token schema is used for authentication codes, access tokens, and
// refresh tokens.
type TokenModel struct {
	ID        string `bson:"_id,omitempty"`
	ClientID  string `bson:"client_id"`
	UserID    string `bson:"user_id"`
	CreatedAt int64  `bson:"createdAt"`
	TTL       int64  `bson:"expireAfterSeconds"`
}

// TokenController defines database operations for the OAuth2 token model.
// It controls refresh and access tokens, as well as authorization codes.
type TokenController interface {

	// Refresh tokens.
	FindRefreshByID(string) (TokenModel, error)
	CreateRefresh(TokenModel) (string, error)
	DeleteRefreshByID(string) error

	// Access tokens.
	FindAccessByID(string) (TokenModel, error)
	CreateAccess(TokenModel) (string, error)
	DeleteAccessByID(string) error

	// Authorization tokens/codes.
	FindAuthByID(string) (TokenModel, error)
	CreateAuth(TokenModel) (string, error)
	DeleteAuthByID(string) error
}

// MongoTokenController implements TokenController using MongoDB.
type MongoTokenController struct {
	state       *MongoState
	refreshColl *mongo.Collection
	accessColl  *mongo.Collection
	authColl    *mongo.Collection
}

// NewTokenController creates a MongoDB user controller using the provided
// database state.
func NewTokenController(state *MongoState) (TokenController, error) {
	if state == nil {
		return nil, ErrNilState
	}

	if state.Stopped.Load() {
		return nil, ErrClosed
	}

	ctrl := new(MongoTokenController)
	ctrl.refreshColl = state.Client.Database(state.Name).Collection("refresh_tokens")
	ctrl.accessColl = state.Client.Database(state.Name).Collection("access_tokens")
	ctrl.authColl = state.Client.Database(state.Name).Collection("auth_tokens")
	ctrl.state = state

	return ctrl, nil
}

func (cc *MongoTokenController) FindRefreshByID(id string) (TokenModel, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.refreshColl == nil {
		return TokenModel{}, ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Extract primitive object ID.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return TokenModel{}, err
	}

	// Find model.
	var token TokenModel
	err = cc.refreshColl.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objID}}).Decode(&token)
	if err != nil {
		return TokenModel{}, err
	}

	return token, nil
}

func (cc *MongoTokenController) CreateRefresh(tk TokenModel) (string, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.refreshColl == nil {
		return "", ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Insert.
	res, err := cc.refreshColl.InsertOne(context.TODO(), tk)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (cc *MongoTokenController) DeleteRefreshByID(id string) error {
	if cc.state == nil || cc.state.Stopped.Load() || cc.refreshColl == nil {
		return ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Extract primitive object ID.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = cc.refreshColl.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: objID}})
	return err
}

func (cc *MongoTokenController) FindAccessByID(id string) (TokenModel, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.accessColl == nil {
		return TokenModel{}, ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Extract primitive object ID.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return TokenModel{}, err
	}

	// Find model.
	var token TokenModel
	err = cc.accessColl.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objID}}).Decode(&token)
	if err != nil {
		return TokenModel{}, err
	}

	return token, nil
}

func (cc *MongoTokenController) CreateAccess(tk TokenModel) (string, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.accessColl == nil {
		return "", ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Insert.
	res, err := cc.accessColl.InsertOne(context.TODO(), tk)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (cc *MongoTokenController) DeleteAccessByID(id string) error {
	if cc.state == nil || cc.state.Stopped.Load() || cc.accessColl == nil {
		return ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Extract primitive object ID.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = cc.accessColl.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: objID}})
	return err
}

func (cc *MongoTokenController) FindAuthByID(id string) (TokenModel, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.authColl == nil {
		return TokenModel{}, ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Extract primitive object ID.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return TokenModel{}, err
	}

	// Find model.
	var token TokenModel
	err = cc.authColl.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objID}}).Decode(&token)
	if err != nil {
		return TokenModel{}, err
	}

	return token, nil
}

func (cc *MongoTokenController) CreateAuth(tk TokenModel) (string, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.authColl == nil {
		return "", ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Insert.
	res, err := cc.authColl.InsertOne(context.TODO(), tk)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (cc *MongoTokenController) DeleteAuthByID(id string) error {
	if cc.state == nil || cc.state.Stopped.Load() || cc.authColl == nil {
		return ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Extract primitive object ID.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = cc.authColl.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: objID}})
	return err
}
