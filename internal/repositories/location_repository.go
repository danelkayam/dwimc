package repositories

import (
	"context"
	"dwimc/internal/model"
	"dwimc/internal/utils"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const COLLECTION_NAME_LOCATIONS = "locations"

type LocationRepository interface {
	GetAllByDevice(deviceID string) ([]model.Location, error)
	GetLatestByDevice(deviceID string) (*model.Location, error)
	Create(deviceID string, latitude float64, longitude float64) (*model.Location, error)
	Delete(deviceID string, id string) (bool, error)
	DeleteAllByDevice(deviceID string) (bool, error)
}

type MongodbLocationRepository struct {
	context    context.Context
	collection *mongo.Collection
}

func NewMongodbLocationRepository(
	context context.Context,
	client *mongo.Client,
	dbName string,
) (LocationRepository, error) {
	collection := client.Database(dbName).Collection(COLLECTION_NAME_LOCATIONS)

	if _, err := collection.Indexes().CreateOne(
		context,
		mongo.IndexModel{
			Keys: bson.M{
				"deviceId": 1,
			},
			Options: options.Index().SetUnique(false),
		}); err != nil {
		return nil, utils.AsError(model.ErrDatabase, err.Error())
	}

	return &MongodbLocationRepository{
		context:    context,
		collection: collection,
	}, nil
}

func (r *MongodbLocationRepository) GetAllByDevice(deviceID string) ([]model.Location, error) {
	objectID, err := bson.ObjectIDFromHex(deviceID)
	if err != nil {
		return nil, utils.AsError(
			model.ErrInvalidArgs,
			fmt.Sprintf("invalid id: %s", deviceID),
		)
	}

	locations := []model.Location{}

	cursor, err := r.collection.Find(
		r.context,
		bson.M{"deviceId": objectID},
	)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return locations, nil
		}
		return nil, err
	}

	defer cursor.Close(r.context)

	for cursor.Next(r.context) {
		var location model.Location

		if err := cursor.Decode(&location); err != nil {
			return nil, utils.AsError(model.ErrDatabase, err.Error())
		}

		locations = append(locations, location)
	}

	return locations, nil
}

func (r *MongodbLocationRepository) GetLatestByDevice(deviceID string) (*model.Location, error) {
	objectID, err := bson.ObjectIDFromHex(deviceID)
	if err != nil {
		return nil, utils.AsError(
			model.ErrInvalidArgs,
			fmt.Sprintf("invalid id: %s", deviceID),
		)
	}

	var location model.Location

	err = r.collection.FindOne(
		r.context,
		bson.M{"deviceId": objectID},
		options.FindOne().SetSort(bson.D{{Key: "updatedAt", Value: -1}}),
	).Decode(&location)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.AsError(model.ErrItemNotFound, "device not found")
		}

		return nil, utils.AsError(model.ErrDatabase, err.Error())
	}

	return &location, nil
}

func (r *MongodbLocationRepository) Create(deviceID string, latitude float64, longitude float64) (*model.Location, error) {
	objectID, err := bson.ObjectIDFromHex(deviceID)
	if err != nil {
		return nil, utils.AsError(
			model.ErrInvalidArgs,
			fmt.Sprintf("invalid id: %s", deviceID),
		)
	}

	created := time.Now().UTC()
	location := &model.Location{
		ID:        bson.NewObjectID(),
		CreatedAt: created,
		UpdatedAt: created,
		DeviceID:  objectID,
		Latitude:  latitude,
		Longitude: longitude,
	}

	result, err := r.collection.InsertOne(r.context, location)
	if err != nil {
		return nil, utils.AsError(model.ErrOperationFailed, err.Error())
	}

	if result.InsertedID == nil {
		return nil, utils.AsError(model.ErrOperationFailed, "failed to insert location")
	}

	return location, nil
}

func (r *MongodbLocationRepository) Delete(deviceID string, id string) (bool, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, utils.AsError(
			model.ErrInvalidArgs,
			fmt.Sprintf("invalid id: %s", id),
		)
	}

	deviceOID, err := bson.ObjectIDFromHex(deviceID)
	if err != nil {
		return false, utils.AsError(
			model.ErrInvalidArgs,
			fmt.Sprintf("invalid id: %s", deviceID),
		)
	}

	result, err := r.collection.DeleteOne(
		r.context,
		bson.M{
			"_id":      objectID,
			"deviceId": deviceOID,
		},
	)

	if err != nil {
		return false, utils.AsError(model.ErrDatabase, err.Error())
	}

	return result.DeletedCount > 0, nil
}

func (r *MongodbLocationRepository) DeleteAllByDevice(deviceID string) (bool, error) {
	objectID, err := bson.ObjectIDFromHex(deviceID)
	if err != nil {
		return false, utils.AsError(
			model.ErrInvalidArgs,
			fmt.Sprintf("invalid id: %s", deviceID),
		)
	}

	result, err := r.collection.DeleteMany(
		r.context,
		bson.M{"deviceId": objectID},
	)

	if err != nil {
		return false, utils.AsError(model.ErrDatabase, err.Error())
	}

	return result.DeletedCount > 0, nil
}
