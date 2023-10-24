package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MongoDB encapsulates a Mongo driver connection to some database.
type MongoDB struct {
	client    *mongo.Client
	isStopped bool
	dbname    string
}

// NewMongoDB creates and initializes a new MongoDB instance with the given
// connection URI. Data will be read/written from dbname.
func NewMongoDB(dbname string, uri string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	m := new(MongoDB)
	m.client = client
	m.isStopped = false
	m.dbname = dbname

	return m, nil
}

// Stop shuts down the mongoDB driver. Once stopped, a MongoDB instance may
// not be restarted.
func (m *MongoDB) Stop() error {
	if m.isStopped {
		return fmt.Errorf("mongoDB driver already stopped")
	}

	if err := m.client.Disconnect(context.TODO()); err != nil {
		return err
	}

	return nil
}
