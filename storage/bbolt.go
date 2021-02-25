package storage

import (
	"go.etcd.io/bbolt"
)

// bboltStore represents a data store using a bbolt file as the underlying storage. Its interface allows for common
// operations to share code.
type bboltStore interface {
	DB() (db *bbolt.DB)
	BucketName() (bucketName []byte)
}

// forEachFunc is a function signature that instructs on what to do with each shortened URL and its data.
type forEachFunc func(shortened, data []byte) (err error)

// bboltDelete deletes the given shortened URLs from the bbolt storage.
func bboltDelete(b bboltStore, shortenedURLs []string) (err error) {

	// Check for the nil case.
	if shortenedURLs == nil || len(shortenedURLs) == 0 {

		// Open the bbolt database for exclusive writing.
		if err = b.DB().Update(func(tx *bbolt.Tx) error {

			// Delete the bucket.
			if err = tx.DeleteBucket(b.BucketName()); err != nil {
				return err
			}

			// Recreate the bucket.
			if _, err = tx.CreateBucket(b.BucketName()); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	} else {

		// Open the bbolt database for writing, batch if possible.
		if err = b.DB().Batch(func(tx *bbolt.Tx) error {

			// Iterate through the given shortened URLs.
			for _, shortened := range shortenedURLs {

				// Delete the shortened URL's Terse data from the bucket.
				if err = tx.Bucket(b.BucketName()).Delete([]byte(shortened)); err != nil {
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

// bboltRead reads the given shortened URLs from the bbolt storage and performs a function on each of their values.
func bboltRead(b bboltStore, forEach forEachFunc, shortenedURLs []string) (err error) {

	// Open the bbolt database for reading.
	if err = b.DB().View(func(tx *bbolt.Tx) error {

		// Check for the nil case.
		if shortenedURLs == nil || len(shortenedURLs) == 0 {

			// Iterate through all the shortened URLs.
			if err = tx.Bucket(b.BucketName()).ForEach(forEach); err != nil {
				return err
			}
		} else {

			// Iterate through the given shortened URLs.
			for _, shortened := range shortenedURLs {

				// Get the raw Terse data.
				data := tx.Bucket(b.BucketName()).Get([]byte(shortened))
				if data == nil {
					return ErrShortenedNotFound
				}

				// Perform the given function on the shortened URL and its data.
				if err = forEach([]byte(shortened), data); err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
