package services

import (
	"context"
	"github.com/rekamarket/mongodb-storage-lib/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaseStorage
type BaseStorage struct {
	Client         *mongo.Client
	DBName         string
	CollectionName string
}

// Model
type Model interface {
	GetID() primitive.ObjectID
	GetHexID() string
	SetHexID(hexID string) error
	SetupTimestamps()
}

// NewBaseStorage() is a constructor for BaseStorage struct
func NewBaseStorage(ctx context.Context, mongoURI, dbName, collectionName string) (*BaseStorage, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))

	if err != nil {
		return nil, err
	}

	return &BaseStorage{
		Client:         client,
		DBName:         dbName,
		CollectionName: collectionName,
	}, nil
}

// Ping() the mongo server
func (bs *BaseStorage) Ping(ctx context.Context) error {
	return bs.Client.Ping(ctx, nil)
}

// GetCollection() returns collections storage
func (bs *BaseStorage) GetCollection() *mongo.Collection {
	return bs.Client.Database(bs.DBName).Collection(bs.CollectionName)
}

// InsertOne() inserts given Model and return an ID of inserted document
func (bs *BaseStorage) InsertOne(ctx context.Context, m Model, opts ...*options.InsertOneOptions) (string, error) {
	m.SetupTimestamps()

	b, err := bson.Marshal(m)

	if err != nil {
		return "", err
	}

	res, err := bs.GetCollection().InsertOne(ctx, b, opts...)

	if err != nil {
		return "", helpers.HandleDuplicationErr(err)
	}

	objectID, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return "", helpers.ErrInvalidObjectID
	}

	hexID := objectID.Hex()

	if err = m.SetHexID(hexID); err != nil {
		return "", nil
	}

	return hexID, nil
}

// InsertMany()
func (bs *BaseStorage) InsertMany(ctx context.Context, docs []interface{}, opts ...*options.InsertManyOptions) ([]string, error) {
	res, err := bs.GetCollection().InsertMany(ctx, docs, opts...)

	if err != nil {
		return []string{}, helpers.HandleDuplicationErr(err)
	}

	hexIDs := []string{}

	for _, insertedID := range res.InsertedIDs {
		objectID, ok := insertedID.(primitive.ObjectID)

		if !ok {
			return []string{}, helpers.ErrInvalidObjectID
		}

		hexIDs = append(hexIDs, objectID.Hex())
	}

	return hexIDs, nil
}

// DropAll() deletes collection from database
func (bs *BaseStorage) DropAll(ctx context.Context) error {
	return bs.GetCollection().Drop(ctx)
}
