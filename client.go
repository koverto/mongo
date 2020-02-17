package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func (c *client) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return c.Database().Collection(name, opts...)
}

func (c *client) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.Client.Connect(ctx); err != nil {
		return err
	}

	pingCtx, cancelPing := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelPing()

	if err := c.Client.Ping(pingCtx, nil); err != nil {
		return err
	}

	if err := c.createIndexes(); err != nil {
		return err
	}

	return nil
}

func (c *client) Database() *mongo.Database {
	return c.Client.Database(c.name)
}

func (c *client) DefineIndexes(indexSets ...IndexSet) {
	c.indexSets = append(c.indexSets, indexSets...)
}

func (c *client) createIndexes() error {
	createOpts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	for _, indexSet := range c.indexSets {
		iv := c.Collection(indexSet.Collection).Indexes()

		if _, err := iv.CreateMany(context.Background(), indexSet.Indexes, createOpts); err != nil {
			return err
		}
	}

	return nil
}
