package repositories

import (
	"context"
	"dwimc/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const COLLECTION_NAME_LOCATIONS = "locations"

type LocationRepository interface {
	GetLocations(serial string) ([]model.Location, error)
	GetLatestLocation(serial string) (*model.Location, error)
	CreateLocation(location model.Location) (*model.Location, error)
	DeleteLocation(id string) error
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
				"deviceSerial": 1,
			},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return nil, err
	}

	return &MongodbLocationRepository{
		context:    context,
		collection: collection,
	}, nil
}

func (r *MongodbLocationRepository) GetLocations(serial string) ([]model.Location, error) {
	cursor, err := r.collection.Find(
		r.context,
		bson.M{"deviceSerial": serial},
	)

	if err != nil {
		// TODO - handle db errors
		return nil, err
	}

	defer cursor.Close(r.context)
	locations := []model.Location{}

	for cursor.Next(r.context) {
		var location model.Location

		if err := cursor.Decode(&location); err != nil {
			// TODO - handle db errors
			return nil, err
		}

		locations = append(locations, location)
	}

	return locations, nil
}

func (r *MongodbLocationRepository) GetLatestLocation(serial string) (*model.Location, error) {
	var location model.Location

	err := r.collection.FindOne(
		r.context,
		bson.M{"deviceSerial": serial},
		options.FindOne().SetSort(bson.D{{Key: "updatedAt", Value: -1}}),
	).Decode(&location)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// TODO - return ItemNotFoundError?
			return nil, nil
		}
		return nil, err
	}

	return &location, nil
}

func (r *MongodbLocationRepository) CreateLocation(location model.Location) (*model.Location, error) {
	created := time.Now().UTC()

	location.ID = bson.NewObjectID().String()
	location.CreatedAt = created
	location.UpdatedAt = created

	result, err := r.collection.InsertOne(r.context, location)
	if err != nil {
		// TODO - handle db errors
		return nil, err
	}

	if result.InsertedID == nil {
		// TODO - handle db errors
		return nil, nil
	}

	return &location, nil
}

func (r *MongodbLocationRepository) DeleteLocation(id string) error {
	_, err := r.collection.DeleteOne(
		r.context,
		bson.M{"_id": id},
	)

	if err != nil {
		// TODO - handle db errors
		return err
	}

	return nil
}
