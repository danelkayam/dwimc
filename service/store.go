package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	Context    context.Context
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

func (store *Store) Upsert(params UpdateParams) (*Operation, error) {
	updatedAt := time.Now().UTC()

	_, err := store.collection.UpdateOne(
		store.Context,
		bson.M{"serial": params.Serial},
		bson.M{
			"$set": bson.M{
				"serial":    params.Serial,
				"name":      params.Name,
				"position":  params.Position,
				"updatedAt": updatedAt,
			},
			"$setOnInsert": bson.M{"createdAt": updatedAt},
		},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		return &Operation{Success: false}, err
	}

	return &Operation{Success: true}, nil
}

func (store *Store) GetOne(serial string) (*Device, error) {
	var device Device

	err := store.collection.FindOne(store.Context, bson.M{"serial": serial}).Decode(&device)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &device, nil
}

// TODO: implement this
func (store *Store) GetAll() ([]Device, error) {
	cursor, err := store.collection.Find(store.Context, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(store.Context)

	devices := make([]Device, 0)

	for cursor.Next(store.Context) {
		var device Device

		if err = cursor.Decode(&device); err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	return devices, nil
}
