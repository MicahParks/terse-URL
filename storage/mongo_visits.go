package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/MicahParks/terse-URL/models"
)

type MongoDBVisits struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoDBVisits(ctx context.Context, databaseName, collectionName string, opts ...*options.ClientOptions) (visitsStore VisitsStore, err error) {

	// Create the Visits store.
	m := &MongoDBVisits{}

	// Create the MongoDB client.
	if m.client, err = mongo.NewClient(opts...); err != nil {
		return nil, err
	}

	// Connect to the MongoDB server.
	if err = m.client.Connect(ctx); err != nil {
		return nil, err
	}

	// Assign the collection.
	m.coll = m.client.Database(databaseName).Collection(collectionName)

	return m, nil
}

func (m *MongoDBVisits) Close(ctx context.Context) (err error) {

	// Disconnect from MongoDBVisits.
	return m.client.Disconnect(ctx)
}

func (m *MongoDBVisits) DeleteVisits(ctx context.Context, shortened string) (err error) {

	// Create a filter that will be used to match all documents for this shortened URL.
	filter := bson.D{{"_id", shortened}}

	// Delete the visits for the shortened URL.
	if _, err = m.coll.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (m *MongoDBVisits) ReadVisits(ctx context.Context, shortened string) (visits []*models.Visit, err error) {

	// Create a filter that will be used to match all documents for this shortened URL.
	filter := bson.D{{"_id", shortened}}

	// Create a cursor for the collection with the filter.
	var cursor *mongo.Cursor
	if cursor, err = m.coll.Find(ctx, filter); err != nil {
		return nil, err
	}

	// Get all visits matching the filter.
	if err = cursor.All(ctx, visits); err != nil {
		return nil, err
	}

	return visits, nil
}

func (m *MongoDBVisits) AddVisit(ctx context.Context, shortened string, visit *models.Visit) (err error) {

	// Create a filter that will be used to match all documents for this shortened URL.
	filter := bson.D{{"_id", shortened}}

	// Create some BSON that will tell MongoDB to append to the array of visits for this shortened URL.
	push := bson.D{{
		"$push", bson.D{{
			"visits", visit,
		}},
	}}

	// Update the shortened URL's visits.
	opts := options.Update().SetUpsert(true)
	if _, err = m.coll.UpdateOne(ctx, filter, push, opts); err != nil {
		return err
	}

	return nil
}
