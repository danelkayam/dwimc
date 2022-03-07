package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	Context        context.Context
	dbClient   *mongo.Client
	collection *mongo.Collection
}

const COLLECTION_NAME_DEVICES = "devices"

func (store *Store) Init(dbUri string, dbName string) error {
	dbClient, err := mongo.Connect(store.Context, options.Client().ApplyURI(dbUri))

	if err != nil {
		return err
	}

	collection := dbClient.Database(dbName).Collection(COLLECTION_NAME_DEVICES)

	if _, err := collection.Indexes().CreateOne(store.Context, mongo.IndexModel{
		Keys: bson.M{
			"serial": 1,
		},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return err
	}

	store.dbClient = dbClient
	store.collection = collection

	return nil
}

func (store *Store) Close(cctx context.Context) error {
	if store.dbClient != nil {
		if err := store.dbClient.Disconnect(cctx); err != nil {
			return err
		}

		store.dbClient = nil
		store.collection = nil
	}

	return nil
}

// TODO: implement this
func (store *Store) Upsert(params UpdateParams) (*Device, error) {
	return &Device{
		Serial:   params.Serial,
		Name:     params.Name,
		Position: params.Position,
	}, nil
}

// TODO: implement this
func (store *Store) GetOne(image string) (*Device, error) {
	return &Device{
		Serial: "car-serial",
		Name:   "my car",
		Position: Position{
			Latitude:  32.1786076,
			Longitude: 34.9172212,
		},
	}, nil
}

// TODO: implement this
func (store *Store) GetAll() ([]Device, error) {
	return []Device{}, nil
}
