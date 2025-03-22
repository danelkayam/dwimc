package repositories

import "go.mongodb.org/mongo-driver/v2/mongo"

type LocationRepository interface {
	// TODO - implement this
}

type MongodbLocationRepository struct {
	client *mongo.Client
	dbName string
}

func NewMongodbLocationRepository(client *mongo.Client, dbName string) LocationRepository {
	// TODO - set indexes
	return &MongodbLocationRepository{
		client: client, 
		dbName: dbName,
	}
}
