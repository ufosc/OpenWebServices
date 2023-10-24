package db

// Database serves collections of data.
type Database interface {

	// User CRUD interface.
	CreateUser(u User) (string, error)
	ReadUser(id string) (User, error)
	UpdateUser(id string, u User) (User, error)
	DeleteUser(id string) (User, error)

	// Verif CRUD interface.
	CreateVerif(v Verif) (string, error)
	ReadVerif(id string) (Verif, error)
	UpdateVerif(id string, v Verif) (Verif, error)
	DeleteVerif(id string) (Verif, error)

	// Client CRUD interface.
	CreateClient(c Client) (string, error)
	ReadClient(id string) (Client, error)
	UpdateClient(id string, c Client) (Client, error)
	DeleteClient(id string) (Client, error)

	// Grant CRUD interface.
	CreateGrant(g Grant) (string, error)
	ReadGrant(id string) (Grant, error)
	UpdateGrant(id string, g Grant) (Grant, error)
	DeleteGrant(id string) (Grant, error)

	// AToken CRUD interface.
	CreateAToken(tk AToken) (string, error)
	ReadAToken(id string) (AToken, error)
	UpdateAToken(id string, tk AToken) (AToken, error)
	DeleteAToken(id string) (AToken, error)

	// RToken CRUD interface.
	CreateRToken(tk RToken) (string, error)
	ReadRToken(id string) (RToken, error)
	UpdateRToken(id string, tk RToken) (RToken, error)
	DeleteRToken(id string) (RToken, error)
}
