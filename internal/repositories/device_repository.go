package repositories

import (
	"context"
	"dwimc/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const COLLECTION_NAME_DEVICES = "devices"

type DeviceRepository interface {
	GetDevices() ([]model.Device, error)
	GetDevice(serial string) (*model.Device, error)
	CreateDevice(serial, name string) (*model.Device, error)
	DeleteDevice(serial string) error
}

type MongodbDeviceRepository struct {
	context context.Context
	collection *mongo.Collection
}

func NewMongodbDeviceRepository(
	context context.Context,
	client *mongo.Client,
	dbName string,
) (DeviceRepository, error) {
	collection := client.Database(dbName).Collection(COLLECTION_NAME_DEVICES)

	if _, err := collection.Indexes().CreateOne(
		context,
		mongo.IndexModel{
			Keys: bson.M{
				"serial": 1,
			},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return nil, err
	}

	return &MongodbDeviceRepository{
		context: context,
		collection: collection,
	}, nil
}

func (r *MongodbDeviceRepository) GetDevices() ([]model.Device, error) {
	// TODO - implement this
	return nil, nil
}

func (r *MongodbDeviceRepository) GetDevice(serial string) (*model.Device, error) {
	// TODO - implement this
	return nil, nil
}

func (r *MongodbDeviceRepository) CreateDevice(serial, name string) (*model.Device, error) {
	// TODO - implement this
	return nil, nil
}

func (r *MongodbDeviceRepository) DeleteDevice(serial string) error {
	// TODO - implement this
	return nil
}
