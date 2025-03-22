package repositories

import "go.mongodb.org/mongo-driver/v2/mongo"

type DeviceRepository interface {
	// TODO - implement this
}

type MongodbDeviceRepository struct {
	client *mongo.Client
	dbName string
}

func NewMongodbDeviceRepository(client *mongo.Client, dbName string) DeviceRepository {
	// TODO - set indexes
	return &MongodbDeviceRepository{
		client: client,
		dbName: dbName,
	}
}
