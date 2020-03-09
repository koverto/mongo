package mongo

import "go.mongodb.org/mongo-driver/mongo"

// IndexSet defines a set of indexes to create on a MongoDB collection.
type IndexSet struct {
	Collection string
	Indexes    []mongo.IndexModel
}

// NewIndexSet creates a new IndexSet for the given collection containing the
// given indexes.
func NewIndexSet(collection string, indexes ...mongo.IndexModel) IndexSet {
	return IndexSet{
		Collection: collection,
		Indexes:    indexes,
	}
}
