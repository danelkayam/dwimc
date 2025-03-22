package repositories

import (
	"context"
	"dwimc/internal/model"
	"time"

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
	context    context.Context
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
		context:    context,
		collection: collection,
	}, nil
}

func (r *MongodbDeviceRepository) GetDevices() ([]model.Device, error) {
	cursor, err := r.collection.Find(r.context, bson.M{})
	if err != nil {
		// TODO - handle db errors
		return nil, err
	}

	defer cursor.Close(r.context)
	devices := []model.Device{}

	for cursor.Next(r.context) {
		var device model.Device

		if err := cursor.Decode(&device); err != nil {
			// TODO - handle db errors
			return nil, err
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func (r *MongodbDeviceRepository) GetDevice(serial string) (*model.Device, error) {
	var device model.Device

	err := r.collection.FindOne(
		r.context,
		bson.M{"serial": serial},
	).Decode(&device)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// TODO - return ItemNotFoundError?
			return nil, nil
		}

		return nil, err
	}

	return &device, nil
}

func (r *MongodbDeviceRepository) CreateDevice(serial, name string) (*model.Device, error) {
	var device model.Device

	updatedAt := time.Now().UTC()
	filter := bson.M{"serial": serial}
	update := bson.M{
		"$set": bson.M{
			"serial":    serial,
			"name":      name,
			"updatedAt": updatedAt,
		},
		"$setOnInsert": bson.M{"createdAt": updatedAt},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	err := r.collection.FindOneAndUpdate(
		r.context,
		filter,
		update,
		opts,
	).Decode(&device)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// TODO - return ItemNotFoundError
			return nil, nil
		}

		return nil, err
	}

	return &device, nil
}

func (r *MongodbDeviceRepository) DeleteDevice(serial string) error {
	_, err := r.collection.DeleteOne(
		r.context,
		bson.M{"serial": serial},
	)

	if err != nil {
		// TODO - handle db errors
		return err
	}

	return nil
}
