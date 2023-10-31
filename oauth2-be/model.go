package main

import "time"

type UserModel struct {
	ID           string `bson:"_id,omitempty"`
	Email        string `bson:"email"`
	Password     string `bson:"password"`
	FirstName    string `bson:"first_name"`
	LastName     string `bson:"last_name"`
	LastVerified int64  `bson:"last_verified"`
	Created      int64  `bson:"created"`
}

type PendingUserModel struct {
	ID   string        `bson:"_id,omitempty"`
	User UserModel     `bson:"user"`
	TTL  time.Duration `bson:"expireAfterSeconds"`
}

type VerifyEmailModel struct {
	ID    string `bson:"_id,omitempty"`
	Email string `bson:"email"`
	TTL   int64  `bson:"expireAfterSeconds"`
}
