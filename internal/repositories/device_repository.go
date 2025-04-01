package repositories

import (
	"context"
	"dwimc/internal/model"
	"dwimc/internal/utils"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const COLLECTION_NAME_DEVICES = "devices"

type DeviceRepository interface {
	GetAll() ([]model.Device, error)
	Get(id string) (*model.Device, error)
	Exists(id string) (bool, error)
	Create(serial string, name string) (*model.Device, error)
	Delete(id string) (bool, error)
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
		return nil, utils.AsError(model.ErrDatabase, err.Error())
	}

	return &MongodbDeviceRepository{
		context:    context,
		collection: collection,
	}, nil
}

func (r *MongodbDeviceRepository) GetAll() ([]model.Device, error) {
	devices := []model.Device{}

	cursor, err := r.collection.Find(r.context, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return devices, nil
		}
		return nil, err
	}

	defer cursor.Close(r.context)

	for cursor.Next(r.context) {
		var device model.Device

		if err := cursor.Decode(&device); err != nil {
			return nil, utils.AsError(model.ErrDatabase, err.Error())
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func (r *MongodbDeviceRepository) Get(id string) (*model.Device, error) {
	var device model.Device

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, utils.AsError(
			model.ErrInvalidArgs,
			fmt.Sprintf("invalid id: %s", id),
		)
	}

	err = r.collection.FindOne(
		r.context,
		bson.M{"_id": objectID},
	).Decode(&device)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.AsError(model.ErrItemNotFound, "device not found")
		}

		return nil, utils.AsError(model.ErrDatabase, err.Error())
	}

	return &device, nil
}

func (r *MongodbDeviceRepository) Exists(id string) (bool, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, utils.AsError(
			model.ErrInvalidArgs,
			fmt.Sprintf("invalid id: %s", id),
		)
	}

	err = r.collection.FindOne(
		r.context,
		bson.M{"_id": objectID},
	).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, utils.AsError(model.ErrItemNotFound, "device not found")
		}

		return false, utils.AsError(model.ErrDatabase, err.Error())
	}

	return true, nil
}

func (r *MongodbDeviceRepository) Create(serial string, name string) (*model.Device, error) {
	if len(serial) == 0 && len(name) == 0 {
		return nil, utils.AsError(model.ErrInvalidArgs, "Fields are empty")
	}

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
		// we don't handle conflict since we are using upsert
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.AsError(model.ErrItemNotFound, "device not found")
		}

		if errors.Is(err, mongo.ErrNilValue) {
			return nil, utils.AsError(model.ErrInvalidArgs, err.Error())
		}

		return nil, utils.AsError(model.ErrDatabase, err.Error())
	}

	return &device, nil
}

func (r *MongodbDeviceRepository) Delete(id string) (bool, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, utils.AsError(
			model.ErrInvalidArgs,
			fmt.Sprintf("invalid id: %s", id),
		)
	}

	result, err := r.collection.DeleteOne(
		r.context,
		bson.M{"_id": objectID},
	)

	if err != nil {
		return false, utils.AsError(model.ErrDatabase, err.Error())
	}

	return result.DeletedCount > 0, nil
}
