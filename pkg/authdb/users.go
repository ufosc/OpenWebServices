package authdb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserModel is the user schema.
type UserModel struct {
	ID        string   `bson:"_id,omitempty"`
	Email     string   `bson:"email"`
	Password  string   `bson:"password"`
	FirstName string   `bson:"first_name"`
	LastName  string   `bson:"last_name"`
	Realms    []string `bson:"realms"`
	CreatedAt int64    `bson:"createdAt"`
}

// PendingUserModel is a sign up request that is awaiting email
// verification.
type PendingUserModel struct {
	_id       string    `bson:"_id,omitempty"`
	ID        string    `bson:"ID"`
	Email     string    `bson:"email"`
	User      UserModel `bson:"user"`
	TTL       int64     `bson:"expireAfterSeconds"`
	CreatedAt int64     `bson:"createdAt"`
}

// UserController defines database operations for the OAuth2 user model.
type UserController interface {

	// Registered/confirmed users.
	FindByEmail(string) (UserModel, error)
	FindByID(string) (UserModel, error)
	Update(UserModel) (int64, error)
	Create(UserModel) (string, error)
	DeleteByID(string) error
	Batch(n, skip int64) ([]UserModel, error)
	Count() (int64, error)

	// Pending users.
	FindPendingByID(string) (PendingUserModel, error)
	FindPendingByEmail(string) (PendingUserModel, error)
	CreatePending(PendingUserModel) (string, error)
	DeletePendingByID(string) error
}

// MongoUserController implements UserController on MongoDB.
type MongoUserController struct {
	state *MongoState
	pcoll *mongo.Collection
	ccoll *mongo.Collection
}

// NewUserController creates a MongoDB user controller using the provided
// database state.
func NewUserController(state *MongoState) (UserController, error) {
	if state == nil {
		return nil, ErrNilState
	}

	if state.Stopped.Load() {
		return nil, ErrClosed
	}

	ctrl := new(MongoUserController)
	ctrl.pcoll = state.Client.Database(state.Name).Collection("pending_users")
	ctrl.ccoll = state.Client.Database(state.Name).Collection("users")
	ctrl.state = state

	return ctrl, nil
}

// Users.

// FindByEmail returns a user by their email address or error on failure.
func (cc *MongoUserController) FindByEmail(email string) (UserModel, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.ccoll == nil {
		return UserModel{}, ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Find model.
	var user UserModel
	err := cc.ccoll.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&user)
	if err != nil {
		return UserModel{}, err
	}

	return user, nil
}

// FindByID finds a registered user by their given ID.
func (cc *MongoUserController) FindByID(id string) (UserModel, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.ccoll == nil {
		return UserModel{}, ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Extract primitive object ID.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return UserModel{}, err
	}

	// Find model.
	var user UserModel
	err = cc.ccoll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objID}}).Decode(&user)
	if err != nil {
		return UserModel{}, err
	}

	return user, nil
}

// Update the given user and synchronize its state with the database.
func (cc *MongoUserController) Update(usr UserModel) (int64, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.ccoll == nil {
		return 0, ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	objID, err := primitive.ObjectIDFromHex(usr.ID)
	if err != nil {
		return 0, err
	}
	usr.ID = ""

	res, err := cc.ccoll.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: objID}},
		bson.D{{Key: "$set", Value: usr}})

	if err != nil {
		return 0, nil
	}

	return res.ModifiedCount, nil
}

// Create a user and save them to the database.
func (cc *MongoUserController) Create(usr UserModel) (string, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.ccoll == nil {
		return "", ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Insert.
	res, err := cc.ccoll.InsertOne(context.TODO(), usr)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// DeleteByID deletes a user by their ID.
func (cc *MongoUserController) DeleteByID(id string) error {
	if cc.state == nil || cc.state.Stopped.Load() || cc.ccoll == nil {
		return ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Extract primitive object ID.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = cc.ccoll.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: objID}})
	return err
}

// Page returns users in batches, with each p >= 0 returning the subsequent batch of 20 users.
func (cc *MongoUserController) Batch(n, skip int64) ([]UserModel, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.ccoll == nil {
		return []UserModel{}, ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()
	cursor, err := cc.ccoll.Find(context.TODO(), bson.D{},
		options.Find().SetLimit(n).SetSkip(skip))

	if err != nil {
		return []UserModel{}, err
	}

	result := []UserModel{}
	if err := cursor.All(context.TODO(), &result); err != nil {
		return []UserModel{}, err
	}

	return result, nil
}

// Count returns the number of documents in the collection.
func (cc *MongoUserController) Count() (int64, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.ccoll == nil {
		return -1, ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	count, err := cc.ccoll.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return -1, err
	}

	return count, nil
}

// Pending users.

// FindPendingByID finds a pending user by their given ID.
func (cc *MongoUserController) FindPendingByID(id string) (
	PendingUserModel, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.pcoll == nil {
		return PendingUserModel{}, ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Find model.
	var user PendingUserModel
	err := cc.pcoll.FindOne(context.TODO(),
		bson.D{{Key: "ID", Value: id}}).Decode(&user)

	if err != nil {
		return PendingUserModel{}, err
	}

	return user, nil
}

// FindPendingByEmail finds a pending user by their registered email.
func (cc *MongoUserController) FindPendingByEmail(email string) (
	PendingUserModel, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.pcoll == nil {
		return PendingUserModel{}, ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Find model.
	var user PendingUserModel
	err := cc.pcoll.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&user)
	if err != nil {
		return PendingUserModel{}, err
	}

	return user, nil
}

// Create a pending user and save them to the database.
func (cc *MongoUserController) CreatePending(usr PendingUserModel) (
	string, error) {
	if cc.state == nil || cc.state.Stopped.Load() || cc.pcoll == nil {
		return "", ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	// Insert
	_, err := cc.pcoll.InsertOne(context.TODO(), usr)
	if err != nil {
		return "", err
	}

	return usr.ID, nil
}

// DeletePendingByID deletes the pending user with the given id.
func (cc *MongoUserController) DeletePendingByID(id string) error {
	if cc.state == nil || cc.state.Stopped.Load() || cc.pcoll == nil {
		return ErrClosed
	}

	cc.state.Wg.Add(1)
	defer cc.state.Wg.Done()

	_, err := cc.pcoll.DeleteOne(context.TODO(),
		bson.D{{Key: "ID", Value: id}})

	return err
}
