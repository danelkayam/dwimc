package repositories

import (
	"context"
	"dwimc/internal/model"

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
		context: context,
		collection: collection,
	}, nil
}

func (r *MongodbLocationRepository) GetLocations(serial string) ([]model.Location, error) {
	// TODO - implement this
	return nil, nil
}

func (r *MongodbLocationRepository) GetLatestLocation(serial string) (*model.Location, error) {
	// TODO - implement this
	return nil, nil
}

func (r *MongodbLocationRepository) CreateLocation(location model.Location) (*model.Location, error) {
	// TODO - implement this
	return nil, nil
}

func (r *MongodbLocationRepository) DeleteLocation(id string) error {
	// TODO - implement this
	return nil
}
