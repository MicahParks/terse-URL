package storage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/MicahParks/ctxerrgroup"
	"go.etcd.io/bbolt"

	"github.com/MicahParks/terseurl/models"
)

// BboltTerse is a TerseStore implementation that relies on a bbolt file for the backend storage.
type BboltTerse struct {
	db           *bbolt.DB
	createCtx    ctxCreator
	group        *ctxerrgroup.Group
	summaryStore SummaryStore
	terseBucket  []byte
	visitsStore  VisitsStore
}

// NewBboltTerse creates a new BboltTerse given the required assets.
func NewBboltTerse(db *bbolt.DB, createCtx ctxCreator, group *ctxerrgroup.Group, summaryStore SummaryStore, terseBucket []byte, visitsStore VisitsStore) (terseStore TerseStore) {
	return &BboltTerse{
		db:           db,
		createCtx:    createCtx,
		group:        group,
		terseBucket:  terseBucket,
		summaryStore: summaryStore,
		visitsStore:  visitsStore,
	}
}

// Close closes the connection to the underlying storage. The ctxerrgroup will be killed. This will not close the
// connection to the VisitsStore. This implementation has no network activity and ignores the given context.
func (b *BboltTerse) Close(_ context.Context) (err error) {

	// Kill the worker pool.
	b.group.Kill()

	// Close the bbolt database file.
	return b.db.Close()
}

// Delete deletes data according to the del argument. If the VisitsStore is not nil, then the same method will be
// called for the associated VisitsStore. This implementation has no network activity and ignores the given
// context.
func (b *BboltTerse) Delete(ctx context.Context, del models.Delete) (err error) {

	// Delete Visits data if required.
	if del.Visits == nil || *del.Visits && b.visitsStore != nil {
		if err = b.visitsStore.Delete(ctx, del); err != nil {
			return err
		}
	}

	// Check to make sure if Terse data needs to be deleted.
	if del.Terse == nil || *del.Terse {

		// Open the bbolt database for writing.
		if err = b.db.Update(func(tx *bbolt.Tx) error {

			// Delete the Terse data bucket from the database.
			if err = tx.DeleteBucket(b.terseBucket); err != nil {
				return err
			}

			// Create the bucket again.
			if _, err = tx.CreateBucket(b.terseBucket); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// DeleteOne deletes data according to the del argument for the given shortened URL. No error should be given if
// the shortened URL is not found. If the VisitsStore is not nil, then the same method will be called for the
// associated VisitsStore. This implementation has no network activity and ignores the given context.
func (b *BboltTerse) DeleteOne(ctx context.Context, del models.Delete, shortened string) (err error) {

	// Delete Visits data if required.
	if del.Visits == nil || *del.Visits && b.visitsStore != nil {
		if err = b.visitsStore.DeleteOne(ctx, del, shortened); err != nil {
			return err
		}
	}

	// Check to make sure if Terse data needs to be deleted.
	if del.Terse == nil || *del.Terse {

		// Open the bbolt database for writing.
		if err = b.db.Update(func(tx *bbolt.Tx) error {

			// Delete the Terse from the bucket.
			if err = tx.Bucket(b.terseBucket).Delete([]byte(shortened)); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// Export returns a map of shortened URLs to export data. This implementation has no network activity and ignores the
// given context.
func (b *BboltTerse) Export(ctx context.Context) (export map[string]models.Export, err error) {

	// Create the map for all visits.
	allVisits := make(map[string][]*models.Visit)

	// Only get visits if there is a visits store.
	if b.visitsStore != nil {

		// Get all the visits.
		if allVisits, err = b.visitsStore.Export(ctx); err != nil {
			return nil, err
		}
	}

	// Create the export map.
	export = make(map[string]models.Export)

	// Open the bbolt database for reading.
	if err = b.db.View(func(tx *bbolt.Tx) error {

		// For every key in the bucket, add it to the export.
		if err = tx.Bucket(b.terseBucket).ForEach(func(shortened, data []byte) error {

			// Turn the Terse data into a Terse Go structure.
			var terse *models.Terse
			if terse, err = unmarshalTerseData(data); err != nil {
				return err
			}

			// Get the visits.
			var visits []*models.Visit // TODO Need to use make()?
			if b.visitsStore != nil {
				visits = allVisits[string(shortened)]
			}

			// Add the Terse and Visits export to the export map.
			export[string(shortened)] = models.Export{
				Terse:  terse,
				Visits: visits,
			}

			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return export, nil
}

// ExportOne returns a export of Terse and Visit data for a given shortened URL. The error must be
// storage.ErrShortenedNotFound if the shortened URL is not found. This implementation has no network activity and
// ignores the given context.
func (b *BboltTerse) ExportOne(ctx context.Context, shortened string) (export models.Export, err error) {

	// Get the Terse from the bucket.
	var terse *models.Terse
	if terse, err = b.getTerse(shortened); err != nil {
		return models.Export{}, err
	}

	// Get the visits.
	visits := make([]*models.Visit, 0)
	if b.visitsStore != nil {
		if visits, err = b.visitsStore.ExportOne(ctx, shortened); err != nil {
			return models.Export{}, err
		}
	}

	return models.Export{
		Terse:  terse,
		Visits: visits,
	}, nil
}

// Import imports the given export's data. If del is not nil, data will be deleted accordingly. If del is nil, data
// may be overwritten, but unaffected data will be untouched. If the VisitsStore is not nil, then the same method
// will be called for the associated VisitsStore. This implementation has no network activity and ignores the given
// context.
func (b *BboltTerse) Import(ctx context.Context, del *models.Delete, export map[string]models.Export) (err error) {

	// Check if data needs to be deleted before importing.
	if del != nil {
		if err = b.Delete(ctx, *del); err != nil {
			return err
		}
	}

	// Import the Visits data.
	if b.visitsStore != nil {

		// Import the Visits data. Never pass in a deletion data structure, because that already happened above.
		if err = b.visitsStore.Import(ctx, nil, export); err != nil {
			return err
		}
	}

	// Write every shortened URL's Terse data to the bbolt database.
	for shortened, exp := range export {

		// Open the bbolt database for writing, batch if possible.
		if err = b.db.Batch(func(tx *bbolt.Tx) error {

			// Turn the Terse into JSON bytes.
			var value []byte
			if value, err = json.Marshal(exp.Terse); err != nil {
				return err
			}

			// Write the Terse to the bucket.
			if err = tx.Bucket(b.terseBucket).Put([]byte(shortened), value); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// Insert adds a Terse to the TerseStore. The shortened URL will be active after this. The error will be
// storage.ErrShortenedExists if the shortened URL is already present. This implementation has no network activity and
// ignores the given context.
func (b *BboltTerse) Insert(_ context.Context, terse *models.Terse) (err error) {

	// Determine if the Terse is already present.
	if _, err = b.getTerseData(*terse.ShortenedURL); !errors.Is(err, ErrShortenedNotFound) {
		if err != nil {
			return err
		}
		return ErrShortenedExists
	}
	err = nil

	// Write the Terse to the bucket.
	if err = b.writeTerse(terse); err != nil {
		return err
	}

	return nil
}

// Read retrieves all non-Visit Terse data give its shortened URL. A nil visit may be passed in and the visit should
// not be recorded. The error must be storage.ErrShortenedNotFound if the shortened URL is not found. This
// implementation has no network activity and ignores the given context.
func (b *BboltTerse) Read(_ context.Context, shortened string, visit *models.Visit) (terse *models.Terse, err error) {

	// Track the visit to this shortened URL. Do this in a separate goroutine so the response is faster.
	if visit != nil && b.visitsStore != nil {
		ctx, cancel := b.createCtx()
		go b.group.AddWorkItem(ctx, cancel, func(workCtx context.Context) (err error) {
			return b.visitsStore.Add(workCtx, shortened, visit)
		})
	}

	// Get the Terse from the bucket.
	if terse, err = b.getTerse(shortened); err != nil {
		return nil, err
	}

	return terse, nil
}

// SummaryStore returns the VisitsStore.
func (b *BboltTerse) SummaryStore() SummaryStore {
	return b.summaryStore
}

// Update assumes the Terse already exists. It will override all of its values. The error must be
// storage.ErrShortenedNotFound if the shortened URL is not found. This implementation has no network activity and
// ignores the given context.
func (b *BboltTerse) Update(_ context.Context, terse *models.Terse) (err error) {

	// Determine if the Terse is already present.
	if _, err = b.getTerseData(*terse.ShortenedURL); err != nil {
		return err
	}

	// Write the Terse to the bucket.
	if err = b.writeTerse(terse); err != nil {
		return err
	}

	return nil
}

// Upsert will upsert the Terse into the backend storage. This implementation has no network activity and ignores the
// given context.
func (b *BboltTerse) Upsert(_ context.Context, terse *models.Terse) (err error) {

	// Write the Terse to the bucket.
	if err = b.writeTerse(terse); err != nil {
		return err
	}

	return nil
}

// VisitsStore returns the VisitsStore.
func (b *BboltTerse) VisitsStore() VisitsStore {
	return b.visitsStore
}

// getTerse gets the shortened URL's Terse data from the bbolt database.
func (b *BboltTerse) getTerse(shortened string) (terse *models.Terse, err error) {

	// Get the Terse data from the bucket.
	var data []byte
	if data, err = b.getTerseData(shortened); err != nil {
		return nil, err
	}

	// Turn the data into the Go structure.
	if terse, err = unmarshalTerseData(data); err != nil {
		return nil, err
	}

	return terse, nil
}

// getTerseData get the shortened URL's raw Terse data from the bbolt database.
func (b *BboltTerse) getTerseData(shortened string) (data []byte, err error) {

	// Open the bbolt database for reading.
	if err = b.db.View(func(tx *bbolt.Tx) error {

		// Get the Terse from the bucket.
		data = tx.Bucket(b.terseBucket).Get([]byte(shortened))

		return nil
	}); err != nil {
		return nil, err
	}

	// If the data is nil, then no Terse exists for this shortened URL.
	if data == nil {
		return nil, ErrShortenedNotFound
	}

	return data, nil
}

// unmarshalTerseData turns the JSON data into the Go structure.
func unmarshalTerseData(data []byte) (terse *models.Terse, err error) {
	terse = &models.Terse{}
	if err = json.Unmarshal(data, terse); err != nil {
		return nil, err
	}
	return terse, nil
}

// writeTerse writes the Terse data to the bbolt database.
func (b *BboltTerse) writeTerse(terse *models.Terse) (err error) {

	// Turn the Terse into JSON bytes.
	var value []byte
	if value, err = json.Marshal(terse); err != nil {
		return err
	}

	// Open the bbolt database for writing, batch if possible.
	if err = b.db.Batch(func(tx *bbolt.Tx) error {

		// Write the Terse to the bucket.
		if err = tx.Bucket(b.terseBucket).Put([]byte(*terse.ShortenedURL), value); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
