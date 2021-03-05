package auth

import (
	"context"
	"fmt"
	"sync"
)

// MemAuthorization implements the AuthorizationStore interface. It's contents will not survive a service restart.
type MemAuthorization struct {
	authLookup      map[string]UserData            // Data structure 1.
	shortenedLookup map[string]map[string]struct{} // Data structure 2.
	authMux         sync.RWMutex
	shortenedMux    sync.RWMutex
}

func (m *MemAuthorization) Append(_ context.Context, usersShortened map[string]UserData) (err error) {

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.Lock()
	defer m.authMux.Unlock()
	m.shortenedMux.Lock()
	defer m.shortenedMux.Unlock()

	// Iterate through the users in the given Authorization data.
	var ok bool
	for user, userData := range usersShortened {

		// Iterate through the shortened URLs in the given Authorization data.
		for shortened, a := range userData {

			// Confirm there is Authorization data for the given user.
			_, ok = m.authLookup[user]
			if !ok {
				m.authLookup[user] = UserData{}
			}

			// Update data in structure 1.
			m.authLookup[user][shortened] = a

			// Confirm the shortened URL is in data structure 2.
			_, ok = m.shortenedLookup[shortened]
			if !ok {
				m.shortenedLookup[shortened] = make(map[string]struct{})
			}

			// Update data in structure 2.
			m.shortenedLookup[shortened][user] = struct{}{}
		}
	}

	return nil
}

func (m *MemAuthorization) Close(_ context.Context) (err error) {

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.Lock()
	defer m.authMux.Unlock()
	m.shortenedMux.Lock()
	defer m.shortenedMux.Unlock()

	// Delete all Authorization data.
	m.deleteAll()

	return nil
}

func (m *MemAuthorization) DeleteShortened(_ context.Context, shortenedURLs []string) (err error) {

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.Lock()
	defer m.authMux.Unlock()
	m.shortenedMux.Lock()
	defer m.shortenedMux.Unlock()

	// Check for the empty case.
	if len(shortenedURLs) == 0 {

		// Delete all Authorization data.
		m.deleteAll()
	} else {

		// Iterate through the given shortened URLs.
		usersAffected := make(map[string][]string)
		for _, shortened := range shortenedURLs {

			// Get the affected users.
			users, ok := m.shortenedLookup[shortened]
			if !ok {
				continue
			}

			// Delete the shortened URL from data structure 2.
			delete(m.shortenedLookup, shortened)

			// Create a map of users to affected shortened URLs.
			affected := addValue(users, shortened)

			// Add the users to the set of affected users.
			merge(usersAffected, affected)
		}

		// Iterate through the affected users.
		for user, s := range usersAffected {

			// Confirm the user is present in data structure 1.
			_, ok := m.authLookup[user]
			if !ok {
				return fmt.Errorf("couldn't find user in authorization data structure 1 when it was present in authorization data structure 2: %w", ErrKeyNotFound)
			}

			// Delete the shortened URLs from data structure 1.
			for _, shortened := range s {
				delete(m.authLookup[user], shortened)
			}
		}
	}

	return nil
}

func (m *MemAuthorization) DeleteUsers(_ context.Context, users []string) (err error) {

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.Lock()
	defer m.authMux.Unlock()
	m.shortenedMux.Lock()
	defer m.shortenedMux.Unlock()

	// Check for the empty case.
	if len(users) == 0 {

		// Delete all Authorization data.
		m.deleteAll()
	} else {

		// Iterate through the given users.
		shortenedAffected := make(map[string][]string)
		for _, user := range users {

			// Get the affected shortened URLs.
			shortened, ok := m.authLookup[user]
			if !ok {
				continue
			}

			// Delete the user from data structure 1.
			delete(m.authLookup[user], user)

			// Create a map of shortened URLs to affected users.
			affected := addValue2(shortened, user)

			// Add the shortened URLs to the set of affected shortened URLs.
			merge(shortenedAffected, affected)
		}

		// Iterate through the affected shortened URLs. Delete the users from data structure 2.
		for shortened, u := range shortenedAffected {
			if err = m.deleteFromDataStructure2([]string{shortened}, u); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *MemAuthorization) Overwrite(_ context.Context, usersShortened map[string]UserData) (err error) {

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.authMux.Lock()
	defer m.authMux.Unlock()
	m.shortenedMux.Lock()
	defer m.shortenedMux.Unlock()

	// Iterate through the users in the given Authorization data.
	var ok bool
	for user, userData := range usersShortened {

		// Update data in structure 1.
		m.authLookup[user] = userData

		// Get the existing UserData.
		var oldUserData UserData
		if oldUserData, ok = m.authLookup[user]; ok {

			// Delete users in data structure 2 to that have the shortened URL in the old set, but not the new one.
			for shortened := range oldUserData {
				if _, ok = userData[shortened]; !ok {
					if err = m.deleteFromDataStructure2([]string{shortened}, []string{user}); err != nil {
						return err
					}
				}
			}
		}

		// Iterate through the shortened URLs in the given Authorization data.
		for shortened := range userData {

			// Confirm the shortened URL is in data structure 2.
			_, ok = m.shortenedLookup[shortened]
			if !ok {
				m.shortenedLookup[shortened] = make(map[string]struct{})
			}

			// Update data in structure 2.
			m.shortenedLookup[shortened][user] = struct{}{}
		}
	}

	return nil
}

func (m *MemAuthorization) ReadUsers(_ context.Context, users []string) (usersShortened map[string]UserData, err error) {

	// Create the return map.
	usersShortened = make(map[string]UserData, len(users))

	// Lock data structure 1 for async safe use.
	m.authMux.RLock()
	defer m.authMux.RUnlock()

	// Check for the empty case.
	if len(users) == 0 {

		// Use all Authorization data.
		usersShortened = m.authLookup
	} else {

		// Iterate through the given users.
		for _, u := range users {

			// Get the user's Authorization data.
			userData, ok := m.authLookup[u]
			if !ok {
				// TODO
			}

			// Assign the Authorization data to the return map.
			usersShortened[u] = userData
		}
	}

	return usersShortened, nil
}

func (m *MemAuthorization) ReadShortened(_ context.Context, shortenedURLs []string) (shortenedUsers map[string]ShortenedData, err error) {

	// Create the return map.
	shortenedUsers = make(map[string]ShortenedData, len(shortenedURLs))

	// Lock both data structure 1 and data structure 2 for async safe use.
	m.shortenedMux.RLock()
	defer m.shortenedMux.RUnlock()
	m.authMux.RLock()
	defer m.authMux.RUnlock()

	// Check for the empty case.
	var ok bool
	if len(shortenedURLs) == 0 {

		// Iterate through every shortened URL.
		for shortened, users := range m.shortenedLookup {

			// Add the users' Authorization data to the return map.
			m.addUsers(shortened, shortenedUsers, users)
		}
	} else {

		// Iterate through the given shortened URLs.
		for _, shortened := range shortenedURLs {

			// Get the users associated with the shortened URLs.
			var users map[string]struct{}
			users, ok = m.shortenedLookup[shortened]
			if !ok {
				// TODO
			}

			// Add the users' Authorization data to the return map.
			m.addUsers(shortened, shortenedUsers, users)
		}
	}

	return shortenedUsers, nil
}

// addUsers adds the given users' Authorization data to teh shortenedUsers map. This does no locking and is not async
// safe.
func (m *MemAuthorization) addUsers(shortened string, shortenedUsers map[string]ShortenedData, users map[string]struct{}) {

	// Iterate through the associated users.
	var ok bool
	for user := range users {

		// Confirm the user is in data structure 1.
		_, ok = m.authLookup[user]
		if !ok {
			// TODO
		}

		// Get the Authorization data for the associated user.
		var a Authorization
		a, ok = m.authLookup[user][shortened]
		if !ok {
			// TODO
		}

		// Confirm there is already Authorization data present for this shortened URL.
		_, ok = shortenedUsers[shortened]
		if !ok {
			shortenedUsers[shortened] = make(ShortenedData)
		}

		// Add the Authorization data to the return map.
		shortenedUsers[shortened][user] = a
	}
}

// deleteAll deletes all of the Authorization data. It does not lock, so a luck must be used for async safe usage.
func (m *MemAuthorization) deleteAll() {

	// Reassign the Authorization data so it's take by the garbage collector.
	m.authLookup = make(map[string]UserData)
	m.shortenedLookup = make(map[string]map[string]struct{})
}

// deleteFromDataStructure2 TODO
func (m *MemAuthorization) deleteFromDataStructure2(shortenedURLs, users []string) (err error) {

	// Iterate through the given shortened URLs.
	for _, shortened := range shortenedURLs {

		// Confirm the shortened URL is present in data structure 2.
		_, ok := m.shortenedLookup[shortened]
		if !ok {
			return fmt.Errorf("couldn't find shortened URL in authorization data structure 2 when it was present in authorization data structure 1: %w", ErrKeyNotFound)
		}

		// Delete the users from data structure 2.
		for _, user := range users {
			delete(m.shortenedLookup[shortened], user)
		}
	}

	return nil
}

// addValue takes the given map and adds the string value as the value in the map for every key.
func addValue(given map[string]struct{}, value string) (m map[string][]string) { // This function could benefit from generics.
	m = make(map[string][]string, len(given))
	for key := range given {
		m[key] = append(m[key], value)
	}
	return m
}

// addValue2 takes the given map and adds the string value as the value in the map for every key.
func addValue2(given map[string]Authorization, value string) (m map[string][]string) { // This function could benefit from generics.
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
