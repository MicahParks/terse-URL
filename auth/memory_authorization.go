package auth

import (
	"context"
	"fmt"
	"sync"
)

// MemAuthorization implements the AuthorizationStore interface. It's contents will not survive a service restart.
type MemAuthorization struct {
	authMap  map[string]UserData // Data structure 1.
	indexMap *shortenedIndex     // Data structure 2.
	authMux  sync.RWMutex
}

// TODO
func NewMemAuthorization() (authStore AuthorizationStore) {
	return &MemAuthorization{
		authMap:  make(map[string]UserData),
		indexMap: newShortenedLookup(),
	}
}

// TODO
func (m *MemAuthorization) Append(_ context.Context, usersShortened map[string]UserData) (err error) {

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.Lock()
	defer m.authMux.Unlock()

	// Iterate through the users in the given Authorization data.
	m.indexMap.lock(func() {
		var ok bool
		for user, userData := range usersShortened {

			// Iterate through the shortened URLs in the given Authorization data.
			for shortened, a := range userData {

				// Confirm there is Authorization data for the given user.
				_, ok = m.authMap[user]
				if !ok {
					m.authMap[user] = UserData{}
				}

				// Update data in structure 1.
				m.authMap[user][shortened] = a

				// Update data in structure 2.
				m.indexMap.add(map[string]userSet{shortened: {user: {}}})
			}
		}
	})

	return nil
}

// TODO
func (m *MemAuthorization) Close(_ context.Context) (err error) {

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.Lock()
	defer m.authMux.Unlock()

	// Delete all Authorization data.
	m.indexMap.lock(func() {
		m.deleteAll()
	})

	return nil
}

// TODO
func (m *MemAuthorization) DeleteShortened(_ context.Context, shortenedURLs []string) (err error) {

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.Lock()
	defer m.authMux.Unlock()

	// Check for the empty case.
	if len(shortenedURLs) == 0 {

		// Delete all Authorization data.
		m.deleteAll()
	} else {

		usersAffected := make(map[string][]string)
		m.indexMap.lock(func() {

			// Iterate through the given shortened URLs.
			for _, shortened := range shortenedURLs {

				// Get the affected users.
				users := m.indexMap.read([]string{shortened})
				if users[shortened] == nil {
					continue
				}

				// Delete the shortened URL from data structure 2.
				m.indexMap.delete(map[string]userSet{shortened: nil})

				// Create a map of users to affected shortened URLs.
				affected := addValue(users[shortened], shortened)

				// Add the users to the set of affected users.
				merge(usersAffected, affected)
			}
		})

		// Iterate through the affected users.
		for user, s := range usersAffected {

			// Confirm the user is present in data structure 1.
			_, ok := m.authMap[user]
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

// TODO
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
		m.indexMap.lock(func() {
			m.indexMap.delete(shortenedAffected)
		})
	}

	return nil
}

// TODO
func (m *MemAuthorization) Overwrite(_ context.Context, usersShortened map[string]UserData) (err error) {

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.Lock()
	defer m.authMux.Unlock()
	m.indexMap.lock(func() {
		var ok bool

		// Iterate through the users in the given Authorization data.
		for user, userData := range usersShortened {

			// Update data in structure 1.
			m.authMap[user] = userData

			// Get the existing UserData.
			var oldUserData UserData
			if oldUserData, ok = m.authMap[user]; ok {

				// Delete users in data structure 2 to that have the shortened URL in the old set, but not the new one.
				for shortened := range oldUserData {
					if _, ok = userData[shortened]; !ok {
						m.indexMap.delete(map[string]userSet{shortened: {user: {}}})
					}
				}
			}

			// Create the set of users to add.
			uSet := make(map[string]userSet)
			for shortened := range userData {
				uSet[shortened] = map[string]struct{}{user: {}}
			}

			// Add the new shortened URL to users relationships.
			m.indexMap.add(uSet)
		}
	})

	return nil
}

// TODO
func (m *MemAuthorization) ReadUsers(_ context.Context, users []string) (usersShortened map[string]UserData, err error) {

	// Create the return map.
	usersShortened = make(map[string]UserData, len(users))

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
				// TODO
			}

			// Assign the Authorization data to the return map.
			usersShortened[u] = userData
		}
	}

	return usersShortened, nil
}

// TODO
func (m *MemAuthorization) ReadShortened(_ context.Context, shortenedURLs []string) (shortenedUserSet map[string]ShortenedData, err error) {

	// Create the return map.
	shortenedUserSet = make(map[string]ShortenedData, len(shortenedURLs))

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.RLock()
	defer m.authMux.RUnlock()

	// Check for the empty case.
	if len(shortenedURLs) == 0 {

		// Iterate through every shortened URL.
		for shortened, users := range m.indexMap.read(nil) {

			// Add the users' Authorization data to the return map.
			m.addUsers(shortened, shortenedUserSet, users)
		}
	} else {

		// Iterate through the given shortened URLs.
		for _, shortened := range shortenedURLs {

			// Get the users associated with the shortened URLs.
			users := m.indexMap.read([]string{shortened})

			// Add the users' Authorization data to the return map.
			m.addUsers(shortened, shortenedUserSet, users[shortened])
		}
	}

	return shortenedUserSet, nil
}

// addUsers adds the given users' Authorization data to the shortenedUserSet map. This does no locking and is not async
// safe.
func (m *MemAuthorization) addUsers(shortened string, shortenedUserSet map[string]ShortenedData, users userSet) {

	// Iterate through the associated users.
	var ok bool
	for user := range users {

		// Confirm the user is in data structure 1.
		_, ok = m.authMap[user]
		if !ok {
			// TODO
		}

		// Get the Authorization data for the associated user.
		var a Authorization
		a, ok = m.authMap[user][shortened]
		if !ok {
			// TODO
		}

		// Confirm there is already Authorization data present for this shortened URL.
		_, ok = shortenedUserSet[shortened]
		if !ok {
			shortenedUserSet[shortened] = ShortenedData{}
		}

		// Add the Authorization data to the return map.
		shortenedUserSet[shortened][user] = a
	}
}

// deleteAll deletes all of the Authorization data. It does not lock, so a luck must be used for async safe usage.
func (m *MemAuthorization) deleteAll() {

	// Reassign the Authorization data so it's take by the garbage collector.
	m.authMap = make(map[string]UserData)
	m.indexMap = newShortenedLookup()
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
