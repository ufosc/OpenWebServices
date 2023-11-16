package main

// UserModel is the user schema.
type UserModel struct {
	ID           string   `bson:"_id,omitempty"`
	Email        string   `bson:"email"`
	Password     string   `bson:"password"`
	FirstName    string   `bson:"first_name"`
	LastName     string   `bson:"last_name"`
	Realms       []string `bson:"realms"`
	LastVerified int64    `bson:"last_verified"`
	Created      int64    `bson:"created"`
}

// PendingUserModel is a sign up request that is awaiting email
// verification.
type PendingUserModel struct {
	ID    string    `bson:"_id,omitempty"`
	Email string    `bson:"email"`
	User  UserModel `bson:"user"`
	TTL   int64     `bson:"expireAfterSeconds"`
}

// VerifyEmailSigninModel is a sign in request that is awaiting
// email verification.
type VerifyEmailSigninModel struct {
	ID    string `bson:"_id,omitempty"`
	Email string `bson:"email"`
	TTL   int64  `bson:"expireAfterSeconds"`
}

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
	Created      int64    `bson:"created"`
	TTL          int64    `bson:"expireAfterSeconds"`
}

// The Token schema is used for authentication codes, access tokens, and
// refresh tokens.
type Token struct {
	ID       string `bson:"_id,omitempty"`
	ClientID string `bson:"client_id"`
	UserID   string `bson:"user_id"`
	Created  int64  `bson:"created"`
	TTL      int64  `bson:"expireAfterSeconds"`
}
