package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// MemAuthorization implements the AuthorizationStore interface. It's contents will not survive a service restart.
type MemAuthorization struct {
	authMap    map[string]UserAuth // Data structure 1.
	authMux    sync.RWMutex
	shortIndex *shortenedIndex // Data structure 2.
	logger     *zap.SugaredLogger
}

// NewMemAuthorization creates a new MemAuthorization.
func NewMemAuthorization() (authStore AuthorizationStore) {
	return &MemAuthorization{
		authMap:    make(map[string]UserAuth),
		shortIndex: newShortenedIndex(nil),
	}
}

// Append upserts the given Authorization data. It works in an append first, overwrite second fashion. For example,
// if the existing Authorization data indicates that a user currently has the shortened URLs X and Y, but the given
// Authorization data indicates that the same user has the shortened URLs Y and Z, then the resulting Authorization
// data for that user will include X, Y, and Z. Where Y's Authorization data is the most recently received.
//
// Both data structure 1 and data structure 2 should be updated during this operation.
func (m *MemAuthorization) Append(_ context.Context, usersShortened map[string]UserAuth) (err error) {

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.Lock()
	defer m.authMux.Unlock()

	// Iterate through the users in the given Authorization data.
	m.shortIndex.lock(func() {
		var ok bool
		for user, userData := range usersShortened {

			// Iterate through the shortened URLs in the given Authorization data.
			for shortened, a := range userData {

				// Confirm there is Authorization data for the given user.
				_, ok = m.authMap[user]
				if !ok {
					m.authMap[user] = UserAuth{}
				}

				// Update data in structure 1.
				m.authMap[user][shortened] = a

				// Update data in structure 2.
				m.shortIndex.add(map[string]userSet{shortened: {user: {}}})
			}
		}
	})

	return nil
}

// Close closes the connection to the underlying storage.
func (m *MemAuthorization) Close(_ context.Context) (err error) {

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.Lock()
	defer m.authMux.Unlock()

	// Delete all Authorization data.
	m.shortIndex.lock(func() {
		m.deleteAll()
	})

	return nil
}

// DeleteShortened deletes the Authorization data for the given shortened URLs for all users. If shortenedURLs is
// nil or empty, all Authorization data are deleted. No error should be returned if a shortened URL is not found.
//
// This should first interact with data structure 2. During this interaction the affected users should be noted and
// used to update the underlying data structure, data structure 1, afterwards.
func (m *MemAuthorization) DeleteShortened(_ context.Context, shortenedURLs []string) (err error) {

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.Lock()
	defer m.authMux.Unlock()

	// Check for the empty case.
	if len(shortenedURLs) == 0 {

		// Delete all Authorization data.
		m.deleteAll()
	} else {

		usersAffected := make(map[string][]string) // TODO Change data type to userSet?
		m.shortIndex.lock(func() {

			// Iterate through the given shortened URLs.
			for _, shortened := range shortenedURLs {

				// Get the affected users.
				var users map[string]userSet
				if users, err = m.shortIndex.read([]string{shortened}); err != nil {

					// Check for an acceptable expected error.
					if errors.Is(err, ErrKeyNotFound) {
						err = nil
						continue
					}

					// Return the unexpected error.
					return
				}

				// Delete the shortened URL from data structure 2.
				m.shortIndex.delete(map[string]userSet{shortened: nil})

				// Create a map of users to affected shortened URLs.
				affected := addValue(users[shortened], shortened)

				// Add the users to the set of affected users.
				merge(usersAffected, affected)
			}
		})
		if err != nil {
			return err
		}

		// Iterate through the affected users.
		var ok bool
		for user, s := range usersAffected {

			// Confirm the user is present in data structure 1.
			_, ok = m.authMap[user]
			if !ok {
				return fmt.Errorf("couldn't find user in authorization data structure 1 when it was present in authorization data structure 2: %w", ErrKeyNotFound)
			}

			// Delete the shortened URLs from data structure 1.
			for _, shortened := range s {
				delete(m.authMap[user], shortened)
			}
		}
	}

	return nil
}

// DeleteUsers deletes the Authorization data for the given users. If users is nil or empty, all Authorization data
// are deleted. No error should be returned if a user is not found.
//
// This should first interact directly with the underlying data structure, data structure 1. During this interaction
// the affected shortened URLs should be noted and used to update data structure 2 afterwards.
func (m *MemAuthorization) DeleteUsers(_ context.Context, users []string) (err error) {

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.Lock()
	defer m.authMux.Unlock()

	// Check for the empty case.
	if len(users) == 0 {

		// Delete all Authorization data.
		m.deleteAll()
	} else {

		// Iterate through the given users.
		shortenedAffected := make(map[string]userSet)
		for _, user := range users {

			// Get the affected shortened URLs.
			userData, ok := m.authMap[user]
			if !ok {
				continue
			}

			// Delete the user from data structure 1.
			delete(m.authMap[user], user)

			// Update the map of shortened URLs to affected users.
			for shortened := range userData {
				if _, ok = shortenedAffected[shortened]; !ok {
					shortenedAffected[shortened] = userSet{}
				}
				shortenedAffected[shortened][user] = struct{}{}
			}
		}

		// Delete the users from data structure 2.
		m.shortIndex.lock(func() {
			m.shortIndex.delete(shortenedAffected)
		})
	}

	return nil
}

// Overwrite upserts the given Authorization data. It works in an overwrite only fashion. For example, if the
// existing Authorization data indicates that a user currently has the shortened URLs X and Y, but the given
// Authorization data indicates that the same user has the shortened URLs Y and Z, then the resulting Authorization
// data for that user will only include Y and Z. Where Y's Authorization data is the most recently received.
//
// Both data structure 1 and data structure 2 should be updated during this operation.
func (m *MemAuthorization) Overwrite(_ context.Context, usersShortened map[string]UserAuth) (err error) {

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.Lock()
	defer m.authMux.Unlock()
	m.shortIndex.lock(func() {

		// Iterate through the users in the given Authorization data.
		for user, userData := range usersShortened {

			// Update data in structure 1.
			m.authMap[user] = userData

			// Get the existing UserAuth.
			if oldUserData, ok := m.authMap[user]; ok {

				// Delete users in data structure 2 to that have the shortened URL in the old set, but not the new one.
				for shortened := range oldUserData {
					if _, ok = userData[shortened]; !ok {
						m.shortIndex.delete(map[string]userSet{shortened: {user: {}}})
					}
				}
			}

			// Create the set of users to add.
			uSet := make(map[string]userSet)
			for shortened := range userData {
				uSet[shortened] = map[string]struct{}{user: {}}
			}

			// Add the new shortened URL to users relationships.
			m.shortIndex.add(uSet)
		}
	})

	return nil
}

// ReadUsers exports the Authorization data for the given users. If users is nil or empty, then all users'
// Authorization data is expected.
//
// This should interact directly with the underlying data structure, data structure 1.
func (m *MemAuthorization) ReadUsers(_ context.Context, users []string) (usersShortened map[string]UserAuth, err error) {

	// Create the return map.
	usersShortened = make(map[string]UserAuth, len(users))

	// Lock data structure 1 for async safe use.
	m.authMux.RLock()
	defer m.authMux.RUnlock()

	// Check for the empty case.
	if len(users) == 0 {

		// Use all Authorization data.
		usersShortened = m.authMap
	} else {

		// Iterate through the given users.
		for _, u := range users {

			// Get the user's Authorization data.
			userData, ok := m.authMap[u]
			if !ok {
				return nil, fmt.Errorf("%w: the given user was not found: %s", ErrKeyNotFound, u)
			}

			// Assign the Authorization data to the return map.
			usersShortened[u] = userData
		}
	}

	return usersShortened, nil
}

// ReadShortened exports the Authorization data for the given shortened URLs. If shortenedURLs is nil or empty, then
// all shortened URL Authorization data is expected.
//
// This should first interact with data structure 2 for faster lookups, then gather the Authorization data from data
// structure 1.
func (m *MemAuthorization) ReadShortened(_ context.Context, shortenedURLs []string) (shortenedUserSet map[string]ShortenedAuth, err error) {

	// Create the return map.
	shortenedUserSet = make(map[string]ShortenedAuth, len(shortenedURLs))

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.RLock()
	defer m.authMux.RUnlock()
	m.shortIndex.rlock(func() {

		// Check for the empty case.
		if len(shortenedURLs) == 0 {

			// Get every shortened URL.
			var users map[string]userSet
			if users, err = m.shortIndex.read(nil); err != nil {
				return
			}

			// Iterate through every shortened URL.
			for shortened, user := range users {

				// Add the users' Authorization data to the return map.
				m.addUsers(shortened, shortenedUserSet, user)
			}
		} else {

			// Iterate through the given shortened URLs.
			for _, shortened := range shortenedURLs {

				// Get the users associated with the shortened URLs.
				var users map[string]userSet
				if users, err = m.shortIndex.read([]string{shortened}); err != nil {
					return
				}

				// Add the users' Authorization data to the return map.
				m.addUsers(shortened, shortenedUserSet, users[shortened])
			}
		}
	})
	if err != nil {
		return nil, err
	}

	return shortenedUserSet, nil
}

// addUsers adds the given users' Authorization data to the shortenedUserSet map. This does no locking and is not async
// safe.
func (m *MemAuthorization) addUsers(shortened string, shortenedUserSet map[string]ShortenedAuth, users userSet) {

	// Iterate through the associated users.
	var ok bool
	for user := range users {

		// Confirm the user is in data structure 1.
		_, ok = m.authMap[user]
		var a Authorization
		if ok {

			// Get the Authorization data for the associated user.
			a, ok = m.authMap[user][shortened]

			// Confirm the data in the data structures is behaving as expected.
			if !ok {
				m.logger.Warnw("User marked as associated to a shortened URL in authorization data structure 2, but not authorization data structure 1.",
					"user", user,
					"shortened", shortened,
				)
				continue
			}
		} else {
			m.logger.Warnw("Unable to add user who isn't in authorization data structure 1.",
				"user", user,
				"shortened", shortened,
			)
			continue
		}

		// Confirm there is already Authorization data present for this shortened URL.
		_, ok = shortenedUserSet[shortened]
		if !ok {
			shortenedUserSet[shortened] = ShortenedAuth{}
		}

		// Add the Authorization data to the return map.
		shortenedUserSet[shortened][user] = a
	}
}

// deleteAll deletes all of the Authorization data. It does not lock, so a luck must be used for async safe usage.
func (m *MemAuthorization) deleteAll() {

	// Reassign the Authorization data so it's take by the garbage collector.
	m.authMap = make(map[string]UserAuth)
	m.shortIndex = newShortenedIndex(nil)
}

// addValue takes the given map and adds the string value as the value in the map for every key.
func addValue(given map[string]struct{}, value string) (m map[string][]string) { // This function could benefit from generics.
	m = make(map[string][]string, len(given))
	for key := range given {
		m[key] = append(m[key], value)
	}
	return m
}

// merge merges appends all values for the keys in src into the dst map.
func merge(dst, src map[string][]string) {
	for key, value := range src {
		dst[key] = append(dst[key], value...)
	}
}
