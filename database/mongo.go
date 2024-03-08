package database

import (
	"context"
	"fmt"
	"github.com/Olprog59/infimetrics/commons"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
}

func NewMongo() *Mongo {
	return &Mongo{}
}

func (m *Mongo) Connect() (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?compressors=snappy,zlib,zstd", commons.MONGO_USER, commons.MONGO_PASSWORD, commons.MONGO_HOST, commons.MONGO_PORT)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	m.Client = client
	return client, nil
}

func (m *Mongo) Close() error {
	err := m.Client.Disconnect(context.Background())
	if err != nil {
		return err
	}
	return nil
}
