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
func (db *Database) ReadUserByEmail(email string) (UserModel, error) {
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

func (db *Database) ReadUser(id string) (UserModel, error) {
	if db.stopped {
		return UserModel{}, fmt.Errorf("db has stopped")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return UserModel{}, err
	}

	var user UserModel
	coll := db.client.Database(db.dbname).Collection("users")
	err = coll.FindOne(context.TODO(), bson.D{{"_id", objectId}}).Decode(&user)
	if err != nil {
		return UserModel{}, err
	}

	return user, nil
}

func (db *Database) UpdateUser(user UserModel) (int64, error) {
	if db.stopped {
		return 0, fmt.Errorf("db has been stopped")
	}

	objectId, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return 0, err
	}

	user.ID = ""

	coll := db.client.Database(db.dbname).Collection("users")
	res, err := coll.UpdateOne(context.TODO(), bson.D{{"_id", objectId}},
		bson.D{{"$set", user}})

	if err != nil {
		return 0, err
	}

	return res.ModifiedCount, nil
}

func (db *Database) CreateUser(user UserModel) (string, error) {
	if db.stopped {
		return "", fmt.Errorf("db has stopped")
	}

	coll := db.client.Database(db.dbname).Collection("users")
	res, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (db *Database) ReadPendingUser(id string) (PendingUserModel, error) {
	if db.stopped {
		return PendingUserModel{}, fmt.Errorf("db has been stopped")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return PendingUserModel{}, err
	}

	var user PendingUserModel
	coll := db.client.Database(db.dbname).Collection("pending_users")
	err = coll.FindOne(context.TODO(), bson.D{{"_id", objectId}}).Decode(&user)
	if err != nil {
		return PendingUserModel{}, err
	}

	return user, nil
}

func (db *Database) ReadPendingUserByEmail(email string) (PendingUserModel, error) {
	if db.stopped {
		return PendingUserModel{}, fmt.Errorf("db has been stopped")
	}
	var user PendingUserModel
	coll := db.client.Database(db.dbname).Collection("pending_users")
	err := coll.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&user)
	if err != nil {
		return PendingUserModel{}, err
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

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (db *Database) DeletePendingUser(id string) error {
	if db.stopped {
		return fmt.Errorf("db has stopped")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	coll := db.client.Database(db.dbname).Collection("pending_users")
	_, err = coll.DeleteOne(context.TODO(), bson.D{{"_id", objectId}})

	return err
}

func (db *Database) CreateVerifyEmailSignin(verif VerifyEmailSigninModel) (string, error) {
	if db.stopped {
		return "", fmt.Errorf("db has been stopped")
	}

	coll := db.client.Database(db.dbname).Collection("verify_email_signin")
	res, err := coll.InsertOne(context.TODO(), verif)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (db *Database) ReadVerifyEmailSigninByEmail(email string) (VerifyEmailSigninModel, error) {
	if db.stopped {
		return VerifyEmailSigninModel{}, fmt.Errorf("db has been stopped")
	}
	var verif VerifyEmailSigninModel
	coll := db.client.Database(db.dbname).Collection("verify_email_signin")
	err := coll.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&verif)
	if err != nil {
		return VerifyEmailSigninModel{}, err
	}

	return verif, nil
}

func (db *Database) ReadVerifyEmailSignin(id string) (VerifyEmailSigninModel, error) {
	if db.stopped {
		return VerifyEmailSigninModel{}, fmt.Errorf("db has been stopped")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return VerifyEmailSigninModel{}, err
	}

	var verif VerifyEmailSigninModel
	coll := db.client.Database(db.dbname).Collection("verify_email_signin")
	err = coll.FindOne(context.TODO(), bson.D{{"_id", objectId}}).Decode(&verif)
	if err != nil {
		return VerifyEmailSigninModel{}, err
	}

	return verif, err
}

func (db *Database) DeleteVerifyEmailSignin(id string) error {
	if db.stopped {
		return fmt.Errorf("db has been stopped")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	coll := db.client.Database(db.dbname).Collection("verify_email_signin")
	_, err = coll.DeleteOne(context.TODO(), bson.D{{"_id", objectId}})

	return err
}

func (db *Database) ReadClient(id string) (ClientModel, error) {
	if db.stopped {
		return ClientModel{}, fmt.Errorf("db has stopped")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ClientModel{}, err
	}

	var client ClientModel
	coll := db.client.Database(db.dbname).Collection("clients")
	err = coll.FindOne(context.TODO(), bson.D{{"_id", objectId}}).Decode(&client)

	if err != nil {
		return ClientModel{}, err
	}

	return client, nil
}

func (db *Database) ReadClientByName(name string) (ClientModel, error) {
	if db.stopped {
		return ClientModel{}, fmt.Errorf("db has stopped")
	}

	var client ClientModel
	coll := db.client.Database(db.dbname).Collection("clients")
	err := coll.FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&client)
	if err != nil {
		return ClientModel{}, err
	}

	return client, nil
}

func (db *Database) CreateClient(client ClientModel) (string, error) {
	if db.stopped {
		return "", fmt.Errorf("db has stopped")
	}

	coll := db.client.Database(db.dbname).Collection("clients")
	res, err := coll.InsertOne(context.TODO(), client)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (db *Database) DeleteClient(id string) error {
	if db.stopped {
		return fmt.Errorf("db has been stopped")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	coll := db.client.Database(db.dbname).Collection("clients")
	_, err = coll.DeleteOne(context.TODO(), bson.D{{"_id", objectId}})

	return err
}

func (db *Database) ReadRefreshToken(id string) (Token, error) {
	if db.stopped {
		return Token{}, fmt.Errorf("db has been stopped")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Token{}, err
	}

	var token Token
	coll := db.client.Database(db.dbname).Collection("refresh_tokens")
	err = coll.FindOne(context.TODO(), bson.D{{"_id", objectId}}).Decode(&token)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}

func (db *Database) DeleteRefreshToken(id string) error {
	if db.stopped {
		return fmt.Errorf("db has been stopped")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	coll := db.client.Database(db.dbname).Collection("refresh_tokens")
	_, err = coll.DeleteOne(context.TODO(), bson.D{{"_id", objectId}})

	return err
}

func (db *Database) CreateRefreshToken(tk Token) (string, error) {
	if db.stopped {
		return "", fmt.Errorf("db has stopped")
	}

	coll := db.client.Database(db.dbname).Collection("refresh_tokens")
	res, err := coll.InsertOne(context.TODO(), tk)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (db *Database) ReadAccessToken(id string) (Token, error) {
	if db.stopped {
		return Token{}, fmt.Errorf("db has been stopped")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Token{}, nil
	}

	var token Token
	coll := db.client.Database(db.dbname).Collection("access_tokens")
	err = coll.FindOne(context.TODO(), bson.D{{"_id", objectId}}).Decode(&token)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}

func (db *Database) DeleteAccessToken(id string) error {
	if db.stopped {
		return fmt.Errorf("db has been stopped")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	coll := db.client.Database(db.dbname).Collection("access_tokens")
	_, err = coll.DeleteOne(context.TODO(), bson.D{{"_id", objectId}})

	return err
}

func (db *Database) CreateAccessToken(tk Token) (string, error) {
	if db.stopped {
		return "", fmt.Errorf("db has stopped")
	}

	coll := db.client.Database(db.dbname).Collection("access_tokens")
	res, err := coll.InsertOne(context.TODO(), tk)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (db *Database) ReadAuthCode(id string) (Token, error) {
	if db.stopped {
		return Token{}, fmt.Errorf("db has been stopped")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Token{}, nil
	}

	var token Token
	coll := db.client.Database(db.dbname).Collection("auth_tokens")
	err = coll.FindOne(context.TODO(), bson.D{{"_id", objectId}}).Decode(&token)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}

func (db *Database) CreateAuthCode(tk Token) (string, error) {
	if db.stopped {
		return "", fmt.Errorf("db has stopped")
	}

	coll := db.client.Database(db.dbname).Collection("auth_tokens")
	res, err := coll.InsertOne(context.TODO(), tk)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (db *Database) DeleteAuthCode(id string) error {
	if db.stopped {
		return fmt.Errorf("db has been stopped")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	coll := db.client.Database(db.dbname).Collection("auth_tokens")
	_, err = coll.DeleteOne(context.TODO(), bson.D{{"_id", objectId}})

	return err
}
