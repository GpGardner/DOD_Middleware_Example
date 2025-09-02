package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoConfig holds the configuration for the MongoDB connection.
type MongoConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// MongoClient wraps the mongo.Client and mongo.Database
type MongoClient struct {
	Client *mongo.Client
	DB     *mongo.Database
}

type Database string

// ConnectMongo initializes the MongoDB client and connects to the database
func ConnectMongo(cfg MongoConfig) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return &MongoClient{Client: client}, nil
}

// FindOne executes a find one operation
func (m *MongoClient) FindOne(ctx context.Context, coll string, filter any, result any) error {
	collection := m.DB.Collection(coll)
	return collection.FindOne(ctx, filter).Decode(result)
}

// FindMany executes a find many operation
func (m *MongoClient) FindMany(ctx context.Context, coll string, filter any, results any) error {
	collection := m.DB.Collection(coll)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, results)
}

// InsertOne inserts a single document
func (m *MongoClient) InsertOne(ctx context.Context, coll string, document any) (*mongo.InsertOneResult, error) {
	collection := m.DB.Collection(coll)
	return collection.InsertOne(ctx, document)
}

// InsertMany inserts multiple documents
func (m *MongoClient) InsertMany(ctx context.Context, coll string, documents []any) (*mongo.InsertManyResult, error) {
	collection := m.DB.Collection(coll)
	return collection.InsertMany(ctx, documents)
}

// UpdateOne performs an update on a single document
func (m *MongoClient) UpdateOne(ctx context.Context, coll string, filter any, update any) (*mongo.UpdateResult, error) {
	collection := m.DB.Collection(coll)
	return collection.UpdateOne(ctx, filter, update)
}

// DeleteOne removes a single document
func (m *MongoClient) DeleteOne(ctx context.Context, coll string, filter any) (*mongo.DeleteResult, error) {
	collection := m.DB.Collection(coll)
	return collection.DeleteOne(ctx, filter)
}

// DeleteMany removes multiple documents
func (m *MongoClient) DeleteMany(ctx context.Context, coll string, filter any) (*mongo.DeleteResult, error) {
	collection := m.DB.Collection(coll)
	return collection.DeleteMany(ctx, filter)
}
