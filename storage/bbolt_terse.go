package storage

import (
	"context"

	"go.etcd.io/bbolt"

	"github.com/MicahParks/terseurl/models"
)

// BboltTerse is a TerseStore implementation that relies on a bbolt file for the backend storage.
type BboltTerse struct {
	db     *bbolt.DB
	bucket []byte
}

// NewBboltTerse creates a new BboltTerse given the required assets.
func NewBboltTerse(db *bbolt.DB, terseBucket []byte) (terseStore TerseStore) {
	return BboltTerse{
		db:     db,
		bucket: terseBucket,
	}
}

// BucketName returns the name of the bbolt bucket.
func (b BboltTerse) BucketName() (bucketName []byte) {
	return b.bucket
}

// Close closes the connection to the underlying storage.
func (b BboltTerse) Close(_ context.Context) (err error) {

	// Close the bbolt database file.
	return b.db.Close()
}

// DB returns the bbolt database.
func (b BboltTerse) DB() (db *bbolt.DB) {
	return b.db
}

// Delete deletes the Terse data for the given shortened URLs. If shortenedURLs is nil or empty, all shortened URL
// Terse data are deleted. There should be no error if a shortened URL is not found.
func (b BboltTerse) Delete(_ context.Context, shortenedURLs []string) (err error) {
	return bboltDelete(b, shortenedURLs)
}

// Read returns a map of shortened URLs to Terse data. If shortenedURLs is nil or empty, all shortened URL Terse
// data are expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
func (b BboltTerse) Read(_ context.Context, shortenedURLs []string) (terseData map[string]*models.Terse, err error) {

	// Create the return map.
	terseData = make(map[string]*models.Terse, len(shortenedURLs))

	// Create the forEachFunc.
	var forEach forEachFunc = func(key, value []byte) (err error) {

		// Turn the raw data into Terse data.
		terse, err := bytesToTerse(value)
		if err != nil {
			return err
		}

		// Add the Terse data to the return map.
		terseData[string(key)] = &terse

		return nil
	}

	// Read the Terse data into the return map.
	if err = bboltRead(b, forEach, shortenedURLs); err != nil {
		return nil, err
	}

	return terseData, nil
}

// Summary summarizes the Terse data for the given shortened URLs. If shortenedURLs is nil or empty, then all
// shortened URL Summary data are expected.
func (b BboltTerse) Summary(_ context.Context, shortenedURLs []string) (summaries map[string]*models.TerseSummary, err error) {

	// Create the return map.
	summaries = make(map[string]*models.TerseSummary, len(shortenedURLs))

	// Create the forEachFunc.
	var forEach forEachFunc = func(shortened, data []byte) (err error) {

		// Turn the raw data into Terse data.
		terse, err := bytesToTerse(data)
		if err != nil {
			return err
		}

		// Add the Terse data to the return map.
		summaries[string(shortened)] = summarizeTerse(terse) // TODO Is map locking required?

		return nil
	}

	// Read the Summary data into the return map.
	if err = bboltRead(b, forEach, shortenedURLs); err != nil {
		return nil, err
	}

	return summaries, nil
}

// Write writes the given Terse data according to the given operation. The error must be storage.ErrShortenedExists
// if an Insert operation cannot be performed due to the Terse data already existing. The error must be
// storage.ErrShortenedNotFound if an Update operation cannot be performed due to the Terse data not existing.
func (b BboltTerse) Write(_ context.Context, terseData map[string]*models.Terse, operation WriteOperation) (err error) {

	// Open the bbolt database for writing, batch if possible.
	if err = b.db.Batch(func(tx *bbolt.Tx) error {

		// Iterate through the given shortened URLs.
		for shortened, terse := range terseData {

			// Check to see if the shortened URL is present in the bucket.
			if operation == Insert || operation == Update {
				value := tx.Bucket(b.bucket).Get([]byte(shortened))
				if value != nil && operation == Insert {
					return ErrShortenedExists
				}
				if value == nil && operation == Update {
					return ErrShortenedNotFound
				}
			}

			// Transform the Terse data into bytes.
			data, err := terseToBytes(*terse) // TODO Check for nil?
			if err != nil {
				return err
			}

			// Write the Terse data to the bucket.
			if err = tx.Bucket(b.bucket).Put([]byte(shortened), data); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
