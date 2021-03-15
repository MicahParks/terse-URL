package storage

import (
	"context"
	"errors"

	"go.etcd.io/bbolt"
	"go.uber.org/zap"
)

// BboltAuthorization implements the AuthorizationStore interface.
type BboltAuthorization struct {
	bucket     []byte
	db         *bbolt.DB
	logger     *zap.SugaredLogger
	shortIndex *shortenedIndex
}

// NewBboltAuthorization creates a new BboltAuthorization.
func NewBboltAuthorization(db *bbolt.DB, authorizationBucket []byte, logger *zap.SugaredLogger) (authStore AuthorizationStore, err error) {

	// Create the AuthorizationStore.
	bboltAuthStore := BboltAuthorization{
		bucket: authorizationBucket,
		db:     db,
		logger: logger,
	}

	// Create the asset needed for data structure 2.
	indexMap := make(map[string]userSet)

	// Create a forEachFunc that will populate data structure 2.
	var ok bool
	var forEach forEachFunc = func(key, value []byte) (err error) {

		// The username is the key.
		user := string(key)

		// Transform the value to a map of shortened URLs to Authorization data.
		var userAuth UserAuth
		if userAuth, err = bytesToUserAuth(value); err != nil {
			return err
		}

		// Iterate through the shortened URLs.
		for shortened := range userAuth {

			// Confirm the shortened URL is already in the indexMap.
			if _, ok = indexMap[shortened]; !ok {
				indexMap[shortened] = userSet{}
			}

			// Add the user to the indexMap.
			indexMap[shortened][user] = struct{}{}
		}

		return nil
	}

	if err = bboltRead(bboltAuthStore, forEach, nil); err != nil {
		return nil, err
	}

	// Assign data structure 2.
	bboltAuthStore.shortIndex = newShortenedIndex(indexMap)

	return bboltAuthStore, nil
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

				// Update data structure 2.
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

	// Check for the empty case.
	if len(shortenedURLs) == 0 {

		// Create a map from the slice of shortened URLs to delete.
		deleteShortened := make(map[string]userSet)
		for _, shortened := range shortenedURLs {
			deleteShortened[shortened] = nil
		}

		// Delete all Authorization data.
		if err = bboltDelete(b, nil); err != nil {
			return err
		}

		// Delete from data structure 2.
		b.shortIndex.lock(func() {
			b.shortIndex.delete(deleteShortened)
		})
	} else {

		// Keep track of the affected users.
		var affectedShortened map[string]userSet

		// Find the affected users.
		b.shortIndex.rlock(func() {
			for _, shortened := range shortenedURLs {

				// Read the affected shortened URLs.
				var affected map[string]userSet
				affected, err = b.shortIndex.read([]string{shortened})
				if err != nil {

					// Check for an acceptable error.
					if errors.Is(err, ErrKeyNotFound) {
						err = nil
					}

					return
				}

				// Add the shortened URL to the map of affected shortened URLs.
				affectedShortened[shortened] = affected[shortened]
			}
		})
		if err != nil {
			return err
		}

		// Turn the map of affected shortened URLs into a map of affected users.
		affectedUsers := flipUserSet(affectedShortened) // TODO See if this function can be used elsewhere.

		// Open the bbolt database for writing, batch if possible.
		if err = b.db.Batch(func(tx *bbolt.Tx) error {

			// Iterate through the affected users.
			for user, shortened := range affectedUsers {

				// Get the map of shortened URLs to Authorization data for the user.
				value := tx.Bucket(b.bucket).Get([]byte(user))
				if value == nil {
					b.logger.Warnw("User found in authorization data structure 2, but not authorization data structure 1.",
						"user", user,
						"function", "delete",
					)
					continue
				}

				// Turn the raw data into Authorization data.
				var userAuth UserAuth
				if userAuth, err = bytesToUserAuth(value); err != nil {
					return err
				}

				// Delete the required shortened URLs.
				for short := range shortened {
					delete(userAuth, short)
				}
			}

			return nil
		}); err != nil {
			return err
		}

		b.shortIndex.lock(func() {
			b.shortIndex.delete(affectedShortened)
		})
	}

	return nil
}

// DeleteUsers deletes the Authorization data for the given users. If users is nil or empty, all Authorization data
// are deleted. No error should be returned if a user is not found.
//
// This should first interact directly with the underlying data structure, data structure 1. During this interaction
// the affected shortened URLs should be noted and used to update data structure 2 afterwards.
func (b BboltAuthorization) DeleteUsers(_ context.Context, users []string) (err error) {

	// Keep track of which shortened URLs are affected.
	affectedShortened := make(map[string]userSet, len(users))

	// Create the forEachFunc.
	var forEach forEachFunc = func(key, value []byte) (err error) {
		var ok bool

		// Turn the raw data into a map of shortened URLs to Authorization data.
		var userAuth UserAuth
		if userAuth, err = bytesToUserAuth(value); err != nil {
			return err
		}

		// Add to the map of affected users.
		shortened := string(key)
		for user := range userAuth {

			// Confirm the shortened URL is already present in the map of affected users.
			if _, ok = affectedShortened[shortened]; !ok {
				affectedShortened[shortened] = userSet{}
			}

			// Add the user to the set of affected users for this shortened URL.
			affectedShortened[shortened][user] = struct{}{}
		}

		return nil
	}

	// Read the affected shortened URLs.
	if err = bboltRead(b, forEach, users); err != nil {
		return err
	}

	// Delete from data structure 1.
	if err = bboltDelete(b, users); err != nil {
		return err
	}

	// Delete from data structure 2.
	b.shortIndex.lock(func() {
		b.shortIndex.delete(affectedShortened)
	})

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
			var ok bool

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
						if _, ok = uAuth[shortened]; !ok {
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

				// Update data structure 2.
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
func (b BboltAuthorization) ReadShortened(_ context.Context, shortenedURLs []string) (shortenedUsers map[string]ShortenedAuth, err error) {

	// Create the return map.
	shortenedUsers = make(map[string]ShortenedAuth, len(shortenedUsers))

	// Use data structure 2 to find the users authorized for the shortened URLs.
	var affectedShortened map[string]userSet
	b.shortIndex.rlock(func() {
		affectedShortened, err = b.shortIndex.read(shortenedURLs)
	})
	if err != nil {
		return nil, err
	}

	// Open the bbolt database for reading.
	var ok bool
	if err = b.db.View(func(tx *bbolt.Tx) error {

		// Iterate through the relevant users.
		for user, shortened := range flipUserSet(affectedShortened) {

			// Get the Authorization data.
			value := tx.Bucket(b.bucket).Get([]byte(user))
			if value == nil {
				b.logger.Warnw("User found in authorization data structure 2, but not authorization data structure 1.",
					"user", user,
					"function", "delete",
				)
				continue
			}

			// Turn the raw data into Authorization data.
			var userAuth UserAuth
			if userAuth, err = bytesToUserAuth(value); err != nil {
				return err
			}

			// Iterate through the desired shortened URLs.
			for short := range shortened {

				// Confirm the shortened URL is in the return map.
				if _, ok = shortenedUsers[short]; !ok {
					shortenedUsers[short] = ShortenedAuth{}
				}

				// Add the Authorization data to the return map.
				shortenedUsers[short][user] = userAuth[user]
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return shortenedUsers, nil
}

// flipUserSet turns a map of affected shortened URLs into a map of affected users.
func flipUserSet(affectedShortened map[string]userSet) (affectedUsers map[string]shortenedSet) {

	// Create the return map.
	affectedUsers = make(map[string]shortenedSet)

	// Iterate through the affected shortened URLs.
	var ok bool
	for shortened, users := range affectedShortened {
		for user := range users {

			// Confirm the user is already in the map of affected users.
			if _, ok = affectedUsers[user]; !ok {
				affectedUsers[user] = shortenedSet{}
			}

			// Add the users shortened URLs to the map of affected users.
			affectedUsers[user][shortened] = struct{}{}
		}
	}

	return affectedUsers
}
