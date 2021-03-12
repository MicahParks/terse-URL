package storage

import (
	"fmt"
	"sync"
)

// Authorization represents Authorization data.
//
// The owner has unrestricted access to a shortened URL's data and is the only one allowed to write Summary data and
// Visits data.
//
// The owner is the only one allowed to change other user's Authorization data for a given shortened URL.
//
// TODO Should this be apart of the swagger spec and communicated to the frontend to disable buttons and such?
type Authorization struct {
	Owner       bool `json:"owner"`
	ReadSummary bool `json:"read_summary"`
	ReadVisits  bool `json:"read_visits"`
	ReadTerse   bool `json:"read_terse"`
	WriteTerse  bool `json:"write_terse"`
}

// ShortenedAuth is a mapping of users to Authorization data.
type ShortenedAuth map[string]Authorization

// UserAuth is a mapping of shortened URLs to Authorization data.
type UserAuth map[string]Authorization

// MergeOverwrite merges the given UserAuth's data into the current one. If there is a conflict, the existing data will
// be overwritten in favor of the given data.
func (u UserAuth) MergeOverwrite(uAuth UserAuth) {
	for shortened, authorization := range uAuth {
		u[shortened] = authorization
	}
}

//  shortenedIndex represents an in memory representation of data structure 2
type shortenedIndex struct { // AKA authorization data structure 2.
	indexMap map[string]userSet
	mux      sync.RWMutex
}

// userSet represents a set of unique user identifiers.
type userSet map[string]struct{}

// newShortenedIndex creates a new ShortenedIndex.
func newShortenedIndex(indexMap map[string]userSet) (shortLookup *shortenedIndex) {
	if indexMap == nil {
		indexMap = make(map[string]userSet)
	}
	return &shortenedIndex{
		indexMap: indexMap,
	}
}

// add adds a mapping of shortened URLs to users.
func (s *shortenedIndex) add(shortenedUserSet map[string]userSet) {

	// Iterate through the given shortened URL user sets.
	for shortened, users := range shortenedUserSet {

		// Confirm the shortened URL exists.
		_, ok := s.indexMap[shortened]
		if !ok {
			s.indexMap[shortened] = userSet{}
		}

		// Iterate through the given users and add them to the set.
		for user := range users {
			s.indexMap[shortened][user] = struct{}{}
		}
	}
}

// delete deletes a mapping of shortened URLs to users. If a shortened URL has no users after this transaction, it is
// removed entirely. A nil userSet indicates the shortened URL should be removed entirely.
func (s *shortenedIndex) delete(shortenedUserSet map[string]userSet) {

	// Iterate through the given shortened URLs.
	for shortened, users := range shortenedUserSet {

		// Check if the shortened URL should be deleted.
		if len(users) == 0 || len(users) == len(s.indexMap[shortened]) {
			delete(s.indexMap, shortened)
			continue
		}

		// Confirm the shortened URL is present in data structure 2.
		_, ok := s.indexMap[shortened]
		if !ok {
			continue
		}

		// Delete the users.
		for user := range users {
			delete(s.indexMap[shortened], user)
		}
	}
}

// lock is a helper function that makes it so the user of this data structure does not need to manage async locking
// directly.
func (s *shortenedIndex) lock(f func()) {
	s.mux.Lock()
	defer s.mux.Unlock()
	f()
}

// rlock is a helper function that makes it so the user of this data structure does not need to manage async locking
// directly.
func (s *shortenedIndex) rlock(f func()) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	f()
}

// read reads all the users for the given shortened URLs. If shortenedURLs is nil, then all data will be returned.
func (s *shortenedIndex) read(shortenedURLs []string) (shortenedUserSet map[string]userSet, err error) {

	// Create the return map.
	shortenedUserSet = make(map[string]userSet)

	// Check for the empty case.
	if shortenedURLs == nil {

		// Return all the data.
		shortenedUserSet = s.indexMap
	} else {

		// Iterate through the given shortened URLs.
		for _, shortened := range shortenedURLs {

			// Get the users associated with this shortened URL.
			users, ok := s.indexMap[shortened]
			if !ok {
				return nil, fmt.Errorf("%w: the given shortened URL was not found: %s", ErrKeyNotFound, shortened)
			}

			// Add the users to the return map.
			shortenedUserSet[shortened] = users
		}
	}

	return shortenedUserSet, nil
}
