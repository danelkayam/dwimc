package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const DB_OP_TIMEOUT_DURATION = 5 * time.Second
const DB_POOL_SIZE = 3

func InitializeDatabase(uri string) (*mongo.Client, error) {
	opts := options.Client().
		ApplyURI(uri).
		SetTimeout(DB_OP_TIMEOUT_DURATION).
		SetMaxPoolSize(DB_POOL_SIZE)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
