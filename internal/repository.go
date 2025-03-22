package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DevicesRepository interface {
	init() error
	close(cctx context.Context) error
	get(serial string) (*Device, error)
	getAll() ([]Device, error)
	update(params UpdateParams) (bool, error)
}

type mongoDevicesRepository struct {
	dbUri      string
	dbName     string
	context    context.Context
	dbClient   *mongo.Client
	collection *mongo.Collection
}

const DB_OP_TIMEOUT_DURATION = 5 * time.Second
const COLLECTION_NAME_DEVICES = "devices"

func CreateDevicesRepository(dbUri string, dbName string) DevicesRepository {
	return &mongoDevicesRepository{
		dbUri:   dbUri,
		dbName:  dbName,
		context: context.Background(),
	}
}

func (repository *mongoDevicesRepository) init() error {
	// sets default parameters to db
	clientOptions := options.Client()
	clientOptions.ApplyURI(repository.dbUri)
	clientOptions.SetTimeout(DB_OP_TIMEOUT_DURATION)

	// TODO - set specific low connection pool size

	dbClient, err := mongo.Connect(clientOptions)
	if err != nil {
		return err
	}

	// TODO - ping the mongodb server to check connectivity

	collection := dbClient.Database(repository.dbName).Collection(COLLECTION_NAME_DEVICES)

	if _, err := collection.Indexes().CreateOne(
		repository.context,
		mongo.IndexModel{
			Keys: bson.M{
				"serial": 1,
			},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return err
	}

	repository.dbClient = dbClient
	repository.collection = collection

	return nil
}

func (repository *mongoDevicesRepository) close(cctx context.Context) error {
	if repository.dbClient != nil {
		if err := repository.dbClient.Disconnect(cctx); err != nil {
			return err
		}

		repository.dbClient = nil
		repository.collection = nil
	}

	return nil
}

func (repository *mongoDevicesRepository) get(serial string) (*Device, error) {
	var device Device

	err := repository.collection.FindOne(
		repository.context,
		bson.M{"serial": serial},
	).Decode(&device)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (repository *mongoDevicesRepository) getAll() ([]Device, error) {
	cursor, err := repository.collection.Find(repository.context, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(repository.context)

	devices := make([]Device, 0)

	for cursor.Next(repository.context) {
		var device Device

		if err = cursor.Decode(&device); err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func (repository *mongoDevicesRepository) update(params UpdateParams) (bool, error) {
	updatedAt := time.Now().UTC()

	_, err := repository.collection.UpdateOne(
		repository.context,
		bson.M{"serial": params.Serial},
		bson.M{
			"$set": bson.M{
				"serial":    params.Serial,
				"name":      params.Name,
				"location":  params.Location,
				"updatedAt": updatedAt,
			},
			"$setOnInsert": bson.M{"createdAt": updatedAt},
		},
		options.UpdateOne().SetUpsert(true),
	)

	if err != nil {
		return false, err
	}

	return true, nil
}
