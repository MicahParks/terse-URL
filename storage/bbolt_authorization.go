package storage

import (
	"context"

	"go.etcd.io/bbolt"
)

// BboltAuthorization implements the AuthorizationStore interface.
type BboltAuthorization struct {
	bucket     []byte
	db         *bbolt.DB
	shortIndex *shortenedIndex
}

// NewBboltAuthorization creates a new BboltAuthorization.
func NewBboltAuthorization(db *bbolt.DB, authorizationBucket []byte) (authStore AuthorizationStore) {
	return BboltAuthorization{
		db:     db,
		bucket: authorizationBucket,
	}
}

// Append upserts the given Authorization data. It works in an append first, overwrite second fashion. For example,
// if the existing Authorization data indicates that a user currently has the shortened URLs X and Y, but the given
// Authorization data indicates that the same user has the shortened URLs Y and Z, then the resulting Authorization
// data for that user will include X, Y, and Z. Where Y's Authorization data is the most recently received.
//
// Both data structure 1 and data structure 2 should be updated during this operation.
func (b BboltAuthorization) Append(_ context.Context, usersShortened map[string]UserAuth) (err error) {

	// Open the bbolt database for writing, batch if possible.
	if err = b.db.Batch(func(tx *bbolt.Tx) error {
		b.shortIndex.lock(func() {

			// Iterate through the given users.
			for user, uAuth := range usersShortened {

				// Confirm the key exists in the bucket. Grab its value.
				var userAuth UserAuth
				value := tx.Bucket(b.bucket).Get([]byte(user))
				if value != nil {

					// Transform the raw data to a map of shortened URLs to Authorization data.
					if userAuth, err = bytesToUserAuth(value); err != nil {
						return
					}

					// Write the new data to the map of shortened URLs to Authorization data.
					userAuth.MergeOverwrite(uAuth)
				} else {

					// Use the given map of shortened URLs to Authorization data since this user was not found.
					userAuth = uAuth
				}

				// Transform the map of shortened URLs to Authorization data into bytes.
				if value, err = userAuthToBytes(userAuth); err != nil {
					return
				}

				// Write the transformed data back to the database.
				if err = tx.Bucket(b.bucket).Put([]byte(user), value); err != nil {
					return
				}

				// Create the set of users to add.
				uSet := make(map[string]userSet)
				for shortened := range uAuth {
					uSet[shortened] = map[string]struct{}{user: {}}
				}

				b.shortIndex.add(uSet)
			}
		})

		return err // This error may be assigned in inner function.
	}); err != nil {
		return err
	}

	return nil
}

// BucketName returns the name of the bbolt bucket.
func (b BboltAuthorization) BucketName() (bucketName []byte) {
	return b.bucket
}

// Close closes the connection to the underlying storage.
func (b BboltAuthorization) Close(_ context.Context) (err error) {

	// Close the bbolt database file.
	return b.db.Close()
}

// DB returns the bbolt database.
func (b BboltAuthorization) DB() (db *bbolt.DB) {
	return b.db
}

// DeleteShortened deletes the Authorization data for the given shortened URLs for all users. If shortenedURLs is
// nil or empty, all Authorization data are deleted. No error should be returned if a shortened URL is not found.
//
// This should first interact with data structure 2. During this interaction the affected users should be noted and
// used to update the underlying data structure, data structure 1, afterwards.
func (b BboltAuthorization) DeleteShortened(_ context.Context, shortenedURLs []string) (err error) {
	// TODO
	panic("implement me")
}

// DeleteUsers deletes the Authorization data for the given users. If users is nil or empty, all Authorization data
// are deleted. No error should be returned if a user is not found.
//
// This should first interact directly with the underlying data structure, data structure 1. During this interaction
// the affected shortened URLs should be noted and used to update data structure 2 afterwards.
func (b BboltAuthorization) DeleteUsers(_ context.Context, users []string) (err error) {

	// Delete from data structure 1.
	if err = bboltDelete(b, users); err != nil {
		return err
	}

	// TODO Delete from data structure 2.

	return nil
}

// Overwrite upserts the given Authorization data. It works in an overwrite only fashion. For example, if the
// existing Authorization data indicates that a user currently has the shortened URLs X and Y, but the given
// Authorization data indicates that the same user has the shortened URLs Y and Z, then the resulting Authorization
// data for that user will only include Y and Z. Where Y's Authorization data is the most recently received.
//
// Both data structure 1 and data structure 2 should be updated during this operation.
func (b BboltAuthorization) Overwrite(_ context.Context, usersShortened map[string]UserAuth) (err error) {

	// Open the bbolt database for writing, batch if possible.
	if err = b.db.Batch(func(tx *bbolt.Tx) error {
		b.shortIndex.lock(func() {

			// Iterate through the given users.
			for user, uAuth := range usersShortened {

				// Confirm the key exists in the bucket. Grab its value.
				value := tx.Bucket(b.bucket).Get([]byte(user))
				if value != nil {

					// Transform the raw data to a map of shortened URLs to Authorization data.
					var userAuth UserAuth
					if userAuth, err = bytesToUserAuth(value); err != nil {
						return
					}

					// Delete users in data structure 2 to that have the shortened URL in the old set, but not the new one.
					for shortened := range userAuth {
						if _, ok := uAuth[shortened]; !ok {
							b.shortIndex.delete(map[string]userSet{shortened: {user: {}}})
						}
					}
				}

				// Transform the map of shortened URLs to Authorization data into bytes.
				if value, err = userAuthToBytes(uAuth); err != nil {
					return
				}

				// Write the transformed data back to the database.
				if err = tx.Bucket(b.bucket).Put([]byte(user), value); err != nil {
					return
				}

				// Create the set of users to add.
				uSet := make(map[string]userSet)
				for shortened := range uAuth {
					uSet[shortened] = map[string]struct{}{user: {}}
				}

				b.shortIndex.add(uSet)
			}
		})

		return err // This error may be assigned in inner function.
	}); err != nil {
		return err
	}

	return nil
}

// ReadUsers exports the Authorization data for the given users. If users is nil or empty, then all users'
// Authorization data is expected.
//
// This should interact directly with the underlying data structure, data structure 1.
func (b BboltAuthorization) ReadUsers(_ context.Context, users []string) (usersShortened map[string]UserAuth, err error) {

	// Create the return map.
	usersShortened = make(map[string]UserAuth, len(users))

	// Create the forEachFunc.
	var forEach forEachFunc = func(key, value []byte) (err error) {

		// Turn the raw data into a map of shortened URLs to Authorization data.
		var userAuth UserAuth
		if userAuth, err = bytesToUserAuth(value); err != nil {
			return err
		}

		// Add the map of shortened URLs to Authorization data to the return map.
		usersShortened[string(key)] = userAuth

		return nil
	}

	// Read the Authorization data into the return map.
	if err = bboltRead(b, forEach, users); err != nil {
		return nil, err
	}

	return usersShortened, nil
}

// ReadShortened exports the Authorization data for the given shortened URLs. If shortenedURLs is nil or empty, then
// all shortened URL Authorization data is expected.
//
// This should first interact with data structure 2 for faster lookups, then gather the Authorization data from data
// structure 1.
func (b BboltAuthorization) ReadShortened(_ context.Context, shortenedURLs []string) (shortenedUserSet map[string]ShortenedAuth, err error) {
	panic("implement me")
}
