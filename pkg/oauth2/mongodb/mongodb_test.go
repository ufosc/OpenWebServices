package mongodb

import (
	"github.com/ufosc/OpenWebServices/pkg/oauth2/db"
	"os"
	"testing"
	"time"
)

// Demo test.
func TestMongoDB(t *testing.T) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		// Execute tests using the command:
		// MONGODB_URI=<URI> go test .
		t.Fatal("MONGODB_URI environment variable not set")
	}

	m, err := NewMongoDB("test", uri)
	if err != nil {
		t.Fatal(err)
	}

	id, err := m.CreateUser(db.User{"user@ufl.edu", "1234567@!abc",
		"john", "doe", time.Now().Unix(), 0, 0})

	if err != nil {
		t.Fatal(err)
	}

	if id == "" {
		t.Errorf("got empty id")
	}
}

// TODO: real unit tests.
