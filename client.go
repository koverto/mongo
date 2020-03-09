// Package mongo provides a simplified wrapper for mongo.Client from the
// official MongoDB library.
package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeout = 10 * time.Second

// Client is the wrapper interface for mongo.Client, containing a connection to
// a single database on a server.
type Client interface {
	Collection(string, ...*options.CollectionOptions) *mongo.Collection
	Connect() error
	Database() *mongo.Database
	DefineIndexes(...IndexSet)
}

type client struct {
	name string
	*mongo.Client
	indexSets []IndexSet
}

// NewClient creates a new Client for the MongoDB server at uri and database
// dbName.
func NewClient(uri string, dbName string) (Client, error) {
	c, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &client{
		name:   dbName,
		Client: c,
	}, nil
}

// Collection returns the collection at name with options.CollectionOptions opts
// for the connected database.
func (c *client) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return c.Database().Collection(name, opts...)
}

// Connect opens the connection to the database, waits for the server to respond
// to ping, and creates the indexes.
func (c *client) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := c.Client.Connect(ctx); err != nil {
		return err
	}

	pingCtx, cancelPing := context.WithTimeout(context.Background(), timeout)
	defer cancelPing()

	if err := c.Client.Ping(pingCtx, nil); err != nil {
		return err
	}

	if err := c.createIndexes(); err != nil {
		return err
	}

	return nil
}

// Database returns the database for the client.
func (c *client) Database() *mongo.Database {
	return c.Client.Database(c.name)
}

// DefineIndexes appends the given indexSets to the Client.
func (c *client) DefineIndexes(indexSets ...IndexSet) {
	c.indexSets = append(c.indexSets, indexSets...)
}

func (c *client) createIndexes() error {
	createOpts := options.CreateIndexes().SetMaxTime(timeout)

	for _, indexSet := range c.indexSets {
		iv := c.Collection(indexSet.Collection).Indexes()

		if _, err := iv.CreateMany(context.Background(), indexSet.Indexes, createOpts); err != nil {
			return err
		}
	}

	return nil
}
