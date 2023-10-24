package db

// User represents an OAuth account. Password is hashed (with a salt).
type User struct {
	Username  string `json:"username" bson:"username" binding:"required"`
	Password  string `json:"password" bson:"password" binding:"required"`
	FirstName string `json:"full_name" bson:"full_name" binding:"required"`
	LastName  string `json:"last_name" bson:"last_name" binding:"required"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
	LastLogin int64  `json:"last_login" bson:"last_login"`

	// Timestamp of last time username (email) was verified. Need to verify
	// UFL.edu emails once per semester to ensure user is still a student.
	LastVerified int64 `json:"last_verified" bson:"last_verified"`
}

// Verif represents a pending account verification email.
type Verif struct {
	User      User  `json:"user" bson:"user" binding:"required"`
	CreatedAt int64 `json:"createdAt" bson:"createdAt"`
	TTL       int64 `bson:"expireAfterSeconds"`
}

// Client represents a registered OAuth client application.
type Client struct {
	Name        string   `json:"name" bson:"name" binding:"required"`
	Type        string   `json:"type" bson:"type" binding:"required"`
	Domain      string   `json:"domain" bson:"domain" binding:"required"`
	RedirectURI string   `json:"redirect_uri" bson:"redirect_uri" binding:"required"`
	Scope       []string `json:"scope" bson:"scope" binding:"required"`
	CreatedAt   int64    `json:"createdAt" bson:"createdAt"`
}

// Grant represents a temporary authorization grant.
type Grant struct {
	ClientID    string `bson:"client_id"`
	RedirectURI string `bson:"redirect_uri"`
	Scope       string `bson:"scope"`
	TTL         int64  `bson:"expireAfterSeconds"`
	CreatedAt   int64  `bson:"createdAt"`
}

// AToken represents a temporary access token.
type AToken struct {
	ClientID  string `bson:"client_id"`
	UserID    string `bson:"user_id"`
	Scope     string `bson:"scope"`
	TTL       int64  `bson:"expireAfterSeconds"`
	CreatedAt int64  `bson:"createdAt"`

	// The Grant that issued this token, if applicable.
	GrantID string
}

// RToken represents a refresh token.
type RToken struct {
	ClientID  string `bson:"client_id"`
	UserID    string `bson:"user_id"`
	ATokenID  string `bson:"atoken_id"`
	Scope     string `bson:"scope"`
	TTL       int64  `bson:"expireAfterSeconds"`
	CreatedAt int64  `bson:"createdAt"`

	// The Grant that issued this token, if applicable.
	GrantID string
}
