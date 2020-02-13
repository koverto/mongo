package mongo

import "go.mongodb.org/mongo-driver/mongo"

type IndexSet struct {
	Collection string
	Indexes    []mongo.IndexModel
}

func NewIndexSet(collection string, indexes ...mongo.IndexModel) IndexSet {
	return IndexSet{
		Collection: collection,
		Indexes:    indexes,
	}
}
