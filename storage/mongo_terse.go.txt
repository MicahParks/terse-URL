package storage

import (
	"context"
	"errors"
	"time"

	"github.com/MicahParks/ctxerrgroup"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/MicahParks/terse-URL/models"
)

type MongoDBTerse struct {
	client      *mongo.Client
	coll        *mongo.Collection
	createCtx   ctxCreator
	errChan     chan<- error
	group       *ctxerrgroup.Group
	visitsStore VisitsStore
}

func NewMongoDBTerse(ctx context.Context, createCtx ctxCreator, databaseName, collectionName string, errChan chan<- error, group *ctxerrgroup.Group, visitsStore VisitsStore, opts ...*options.ClientOptions) (terseStore TerseStore, err error) {

	// Create the Terse store.
	m := &MongoDBTerse{
		createCtx:   createCtx,
		errChan:     errChan,
		group:       group,
		visitsStore: visitsStore,
	}

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

func (m *MongoDBTerse) Close(ctx context.Context) (err error) {

	// Kill all worker goroutines.
	m.group.Kill()

	// Disconnect from MongoDBTerse.
	return m.client.Disconnect(ctx)
}

func (m *MongoDBTerse) ScheduleDeletions(ctx context.Context) (err error) {

	// Create a cursor for the Terse collection.
	var cursor *mongo.Cursor
	if cursor, err = m.coll.Find(ctx, bson.D{}); err != nil {
		return err
	}

	// Iterate through the collection, scheduling deletions asynchronously for Terse that have deletion times.
	for cursor.Next(ctx) {

		// Get the Terse from the database.
		terse := &Terse{}
		if err = cursor.Decode(&terse); err != nil {
			return err
		}

		// If the Terse has a deletion time, schedule it asynchronously.
		if terse.DeleteAt != nil {
			go deleteTerseBlocking(m.createCtx, *terse.DeleteAt, m.errChan, terse.Shortened, m)
		}
	}

	return nil
}

func (m *MongoDBTerse) DeleteTerse(ctx context.Context, shortened string) (err error) {

	// Create a filter that will be used to match all documents for this shortened URL.
	filter := bson.D{{"_id", shortened}}

	// Delete documents matching the filter.
	if _, err = m.coll.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (m *MongoDBTerse) GetTerse(ctx context.Context, shortened string, visit *models.Visit, visitCancel context.CancelFunc, visitCtx context.Context) (original string, err error) {

	// Track the visit to this shortened URL. Do this in a separate goroutine so the response is faster.
	if m.visitsStore != nil {
		go m.group.AddWorkItem(visitCtx, visitCancel, func(workCtx context.Context) (err error) {
			return m.visitsStore.AddVisit(workCtx, shortened, visit)
		})
	}

	// Create a filter to find the Terse for the shortened URL.
	filter := bson.D{{"_id", shortened}}

	// Create a Terse data structure to unmarshal into.
	terse := &Terse{}

	// Find the shortened link's terse and decode the returned Terse.
	if err = m.coll.FindOne(ctx, filter).Decode(terse); err != nil {

		// Transform the not found error, if needed.
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = ErrShortenedNotFound
		}

		return "", err
	}

	// Return original URL in the Terse.
	return terse.Original, nil
}

func (m *MongoDBTerse) UpsertTerse(ctx context.Context, deleteAt *time.Time, original, shortened string) (err error) {

	// Create a filter that will be used to match all documents for this shortened URL.
	filter := bson.D{{"_id", shortened}}

	// Create the data structure to upsert with the shortened key.
	terse := Terse{ // TODO Pointer needed?
		DeleteAt:  deleteAt,
		Original:  original,
		Shortened: shortened,
	}

	// Create an option that specifies to insert the Terse if one does not already exist.
	opts := options.Replace().SetUpsert(true)

	// Upsert the Terse.
	if _, err = m.coll.ReplaceOne(ctx, filter, terse, opts); err != nil {
		return err
	}

	// Schedule the Terse for deletion, if a deletion time was given.
	if deleteAt != nil {
		go deleteTerseBlocking(m.createCtx, *deleteAt, m.errChan, shortened, m)
	}

	return nil
}

func (m *MongoDBTerse) VisitsStore() VisitsStore {
	return m.visitsStore
}
