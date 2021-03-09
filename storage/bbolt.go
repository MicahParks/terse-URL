package storage

import (
	"go.etcd.io/bbolt"
)

// bboltStore represents a data store using a bbolt file as the underlying storage. Its interface allows for common
// operations to share code.
type bboltStore interface {
	BucketName() (bucketName []byte)
	DB() (db *bbolt.DB)
}

// forEachFunc is a function signature that instructs on what to do with each shortened URL and its data.
type forEachFunc func(key, value []byte) (err error)

// bboltDelete deletes the given keys from the bbolt storage.
func bboltDelete(b bboltStore, keys []string) (err error) {

	// Check for the empty case.
	if len(keys) == 0 {

		// Open the bbolt database for exclusive writing.
		if err = b.DB().Update(func(tx *bbolt.Tx) error {

			// Delete the bucket.
			if err = tx.DeleteBucket(b.BucketName()); err != nil {
				return err
			}

			// TODO Does this work or does it need to be committed?

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

			// Iterate through the given keys.
			for _, key := range keys {

				// Delete the keys data from the bucket.
				if err = tx.Bucket(b.BucketName()).Delete([]byte(key)); err != nil {
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

// bboltRead reads the given keys from the bbolt storage and performs a function on each of their values.
func bboltRead(b bboltStore, forEach forEachFunc, keys []string) (err error) {

	// Open the bbolt database for reading.
	if err = b.DB().View(func(tx *bbolt.Tx) error {

		// Check for the empty case.
		if len(keys) == 0 {

			// Iterate through all the keys.
			if err = tx.Bucket(b.BucketName()).ForEach(forEach); err != nil {
				return err
			}
		} else {

			// Iterate through the given keys.
			for _, shortened := range keys {

				// Get the raw data.
				data := tx.Bucket(b.BucketName()).Get([]byte(shortened))
				if data == nil {
					return ErrShortenedNotFound
				}

				// Perform the given function on the key and its data.
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
