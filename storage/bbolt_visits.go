package storage

import (
	"context"
	"encoding/json"

	"go.etcd.io/bbolt"

	"github.com/MicahParks/terseurl/models"
)

// BboltVisits if a VisitsStore implementation that relies on a bbolt file for the backend storage.
type BboltVisits struct {
	db           *bbolt.DB
	visitsBucket []byte
}

// NewBboltVisits creates a new NewBboltVisits given the required assets.
func NewBboltVisits(db *bbolt.DB, visitsBucket []byte) (visitsStore VisitsStore) {
	return &BboltVisits{
		db:           db,
		visitsBucket: visitsBucket,
	}
}

// Add adds the visit to the visits store.
func (b *BboltVisits) Add(_ context.Context, shortened string, visit *models.Visit) (err error) {

	// Get the existing visits.
	var visits []*models.Visit
	if visits, err = b.exportShortened(shortened); err != nil {
		return err
	}

	// Add the visits to the existing visits.
	visits = append(visits, visit)

	// Turn the visits into JSON data.
	var data []byte
	if data, err = json.Marshal(visits); err != nil {
		return err
	}

	// Open the bbolt database for writing, batch if possible.
	if err = b.db.Batch(func(tx *bbolt.Tx) error {

		// Put the updated JSON data into the bucket.
		if err = tx.Bucket(b.visitsBucket).Put([]byte(shortened), data); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// Close lets the garbage collector take care of the old Visits data.
func (b *BboltVisits) Close(_ context.Context) (err error) {
	return b.db.Close()
}

// Delete deletes data according to the del argument.
func (b *BboltVisits) Delete(_ context.Context, del models.Delete) (err error) {

	// Confirm Visits data deletion.
	if del.Visits == nil || *del.Visits {

		// Open the bbolt database for writing.
		if err = b.db.Update(func(tx *bbolt.Tx) error {

			// Delete the Visits from the bucket.
			if err = tx.DeleteBucket(b.visitsBucket); err != nil {
				return err
			}

			// Create the bucket again.
			if _, err = tx.CreateBucket(b.visitsBucket); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// DeleteSome deletes data according to the del argument for the shortened URL. No error will be given if the shortened
// URL is not found.
func (b *BboltVisits) DeleteSome(_ context.Context, del models.Delete, shortenedURLs []string) (err error) {

	// Confirm Visits data deletion.
	if del.Visits == nil || *del.Visits {

		// Open the bbolt database for writing, batch if possible.
		if err = b.db.Batch(func(tx *bbolt.Tx) error {

			// Iterate through the shortened URLs.
			for _, shortened := range shortenedURLs {

				// Delete all of this shortened URL's visits from the bucket.
				if err = tx.Bucket(b.visitsBucket).Delete([]byte(shortened)); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// Export exports all exports all visits data.
func (b *BboltVisits) Export(_ context.Context) (allVisits map[string][]*models.Visit, err error) {

	// Create the return map.
	allVisits = make(map[string][]*models.Visit)

	// Open the bbolt database for reading.
	if err = b.db.View(func(tx *bbolt.Tx) error {

		// Iterate through all the keys.
		if err = tx.Bucket(b.visitsBucket).ForEach(func(shortened, value []byte) error {

			// Create the visits.
			visits := make([]*models.Visit, 0)

			// Unmarshal the visit.
			if err = json.Unmarshal(value, &visits); err != nil {
				return err
			}

			// Assign the visits to the map.
			allVisits[string(shortened)] = visits

			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return allVisits, nil
}

// ExportCounts creates a map of shortened URLs to count of Visits.
func (b *BboltVisits) ExportCounts(_ context.Context) (counts map[string]uint, err error) {

	// Create the return map.
	counts = make(map[string]uint)

	// Open the bbolt database for viewing.
	if err = b.db.View(func(tx *bbolt.Tx) error {

		// Iterate through all the keys.
		if err = tx.Bucket(b.visitsBucket).ForEach(func(shortened, value []byte) error {

			// Create the visits.
			visits := make([]*models.Visit, 0)

			// Unmarshal the visit.
			if err = json.Unmarshal(value, &visits); err != nil {
				return err
			}

			// Assign the number of visits to the map.
			counts[string(shortened)] = uint(len(visits))

			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return counts, nil
}

// ExportSome gets all visits to the shortened URL. The error storage.ErrShortenedNotFound will be given if the shortened
// URL is not found.
func (b *BboltVisits) ExportSome(_ context.Context, shortenedURLs []string) (visits map[string][]*models.Visit, err error) {

	// Create the return map.
	visits = make(map[string][]*models.Visit)

	// Open the bbolt database for reading.
	var data []byte
	if err = b.db.View(func(tx *bbolt.Tx) error {

		// Iterate through the shortened URLs.
		for _, shortened := range shortenedURLs {

			// Get the Visits from the bucket.
			data = tx.Bucket(b.visitsBucket).Get([]byte(shortened))

			// Only unmarshal data if there was any.
			var v []*models.Visit
			if data != nil {

				// Turn the JSON data into the Go structure.
				if err = json.Unmarshal(data, &v); err != nil {
					return err // TODO Check for if shortened URL is not found and confirm same for other methods.
				}
			} else {
				// TODO Missing shortened URL?
			}

			// Add the visits to the return map.
			visits[shortened] = v
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return visits, nil
}

// Import imports the given export's data. If del is not nil, data will be deleted accordingly. If del is nil, data
// may be overwritten, but unaffected data will be untouched.
func (b *BboltVisits) Import(ctx context.Context, del *models.Delete, export map[string]models.Export) (err error) {

	// Check if data needs to be deleted before importing.
	if del != nil {
		if err = b.Delete(ctx, *del); err != nil {
			return err
		}
	}

	// Write every shortened URL's Visits data to the bbolt database.
	for shortened, exp := range export {

		// Open the bbolt database for writing, batch if possible.
		if err = b.db.Batch(func(tx *bbolt.Tx) error {

			// Turn the Terse into JSON bytes.
			var value []byte
			if value, err = json.Marshal(exp.Visits); err != nil {
				return err
			}

			// Write the Visits to the bucket.
			if err = tx.Bucket(b.visitsBucket).Put([]byte(shortened), value); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// exportShortened gets all visits for the given shortened URL.
func (b *BboltVisits) exportShortened(shortened string) (visits []*models.Visit, err error) {

	// Open the bbolt database for reading.
	var data []byte
	if err = b.db.View(func(tx *bbolt.Tx) error {

		// Get the Visits from the bucket.
		data = tx.Bucket(b.visitsBucket).Get([]byte(shortened))

		return nil
	}); err != nil {
		return nil, err
	}

	// Only unmarshal data if there was any.
	if data != nil {

		// Turn the JSON data into the Go structure.
		if err = json.Unmarshal(data, &visits); err != nil {
			return nil, err
		}
	}

	return visits, nil
}
