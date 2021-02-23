package storage

import (
	"context"

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
	return BboltVisits{
		db:           db,
		visitsBucket: visitsBucket,
	}
}

// Close closes the connection to the underlying storage.
func (b BboltVisits) Close(_ context.Context) (err error) {

	// Close the bbolt database file.
	return b.db.Close()
}

// BucketName returns the name of the bbolt bucket.
func (b BboltVisits) BucketName() (bucketName []byte) {
	return b.visitsBucket
}

// DB returns the bbolt database.
func (b BboltVisits) DB() (db *bbolt.DB) {
	return b.db
}

// Delete deletes Visits data for the given shortened URLs. If shortenedURLs is nil, then all Visits data are
// deleted. No error should be given if a shortened URL is not found.
func (b BboltVisits) Delete(_ context.Context, shortenedURLs []string) (err error) {
	return bboltDelete(b, shortenedURLs)
}

// Insert inserts the given Visits data. The visits do not need to be unique, so the Visits data should be appended
// to the data structure in storage.
func (b BboltVisits) Insert(_ context.Context, visitsData map[string][]models.Visit) (err error) {

	// Open the bbolt database for writing, batch if possible.
	if err = b.db.Batch(func(tx *bbolt.Tx) error {

		// Iterate through the given shortened URLs.
		for shortened, visits := range visitsData {

			// Get the existing Visits data.
			var existingVisits []models.Visit
			data := tx.Bucket(b.visitsBucket).Get([]byte(shortened))

			// Transform the raw data into Visits data.
			if data != nil {
				if existingVisits, err = bytesToVisits(data); err != nil {
					return err
				}
			}

			// Add the given visits data to the existing Visits data.
			existingVisits = append(existingVisits, visits...)

			// Turn the Visits data back into raw data.
			if data, err = visitsToBytes(existingVisits); err != nil {
				return err
			}

			// Write the Visits data back to the bucket.
			if err = tx.Bucket(b.visitsBucket).Put([]byte(shortened), data); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return err
}

// Read exports the Visits data for the given shortened URLs. If shortenedURLs is nil, then all shortened URL Visits
// data are expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
func (b BboltVisits) Read(_ context.Context, shortenedURLs []string) (visitsData map[string][]models.Visit, err error) {

	// Create the return map.
	visitsData = make(map[string][]models.Visit)

	// Create the forEachFunc.
	var forEach forEachFunc = func(shortened, data []byte) (err error) {

		// Turn the raw data into Visits data.
		visits, err := bytesToVisits(data)
		if err != nil {
			return err
		}

		// Add the Visits data to the return map.
		visitsData[string(shortened)] = visits

		return nil
	}

	// Read the Visits data into the return map.
	if err = bboltRead(b, forEach, shortenedURLs); err != nil {
		return nil, err
	}

	return visitsData, nil
}

// Summary summarizes the Visits data for the given shortened URLs. If shortenedURLs is nil, then all shortened URL
// Summary data are expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
func (b BboltVisits) Summary(_ context.Context, shortenedURLs []string) (summaries map[string]*models.VisitsSummary, err error) {

	// Create the return map.
	summaries = make(map[string]*models.VisitsSummary)

	// Create the forEachFunc.
	var forEach forEachFunc = func(shortened, data []byte) (err error) {

		// Turn the raw data into Visits data.
		visits, err := bytesToVisits(data)
		if err != nil {
			return err
		}

		// Add the Visits data to the return map.
		summaries[string(shortened)] = &models.VisitsSummary{
			VisitCount: uint64(len(visits)),
		}

		return nil
	}

	// Read the Summary data into the return map.
	if err = bboltRead(b, forEach, shortenedURLs); err != nil {
		return nil, err
	}

	return summaries, nil
}
