package storage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/MicahParks/ctxerrgroup"
	"go.etcd.io/bbolt"

	"github.com/MicahParks/terse-URL/models"
)

// BboltTerse is a TerseStore implementation that relies on a bbolt file for the backend storage.
type BboltTerse struct {
	db          *bbolt.DB
	createCtx   ctxCreator
	group       *ctxerrgroup.Group
	terseBucket []byte
	visitsStore VisitsStore
}

// NewBboltTerse creates a new BboltTerse given the required assets.
func NewBboltTerse(db *bbolt.DB, createCtx ctxCreator, group *ctxerrgroup.Group, terseBucket []byte, visitsStore VisitsStore) (terseStore TerseStore) {
	return &BboltTerse{
		db:          db,
		createCtx:   createCtx,
		group:       group,
		terseBucket: terseBucket,
		visitsStore: visitsStore,
	}
}

// Close closes the bbolt database file and kills the ctxerrgroup. This implementation has no network activity and
// ignores the given context.
func (b *BboltTerse) Close(_ context.Context) (err error) {

	// Kill the worker pool.
	b.group.Kill()

	// Close the bbolt database file.
	return b.db.Close()
}

// Delete deletes Terse and Visits data as instructed. This implementation has no network activity and ignores the given
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

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// DeleteOne deletes the Terse and or Visits data for the given shortened URL as instructed. This implementation has
// no network activity and ignores the given context.
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

// Export exports all Terse and Visits data. This implementation has no network activity and partially ignores the
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

// ExportOne exports Terse and Visits data for the given shortened URL. This implementation has no network activity and
// partially ignores the given context.
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

// Import imports the given export. It will delete existing Terse and Visits data as instructed by the del argument, if
// not nil. If the del argument is not present, data will be overwritten and unaffected data will be left alone. This
// implementation has no network activity and ignores the given context.
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

// Insert inserts the Terse into the TerseStore. It will fail if the Terse already exists. This implementation has
// no network activity and ignores the given context.
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

// Read gets the Terse data for the given shortened URL. This implementation has no network activity and ignores
// the given context.
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

// Update updates the Terse into the TerseStore. It will fail if the Terse does not exist. This implementation has
// no network activity and ignores the given context.
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

// Upsert upserts the Terse into the TerseStore. This implementation has no network activity and ignores the given
// context.
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
