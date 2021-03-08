package auth

import (
	"sync"
)

type userSet map[string]struct{}

// TODO
type shortenedIndex struct { // AKA authorization data structure 2.
	indexMap map[string]userSet
	mux      sync.RWMutex
}

func newShortenedLookup() (shortLookup *shortenedIndex) {
	return &shortenedIndex{
		indexMap: make(map[string]userSet),
	}
}

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

func (s *shortenedIndex) delete(shortenedUserSet map[string]userSet) {

	// Iterate through the given shortened URLs.
	for shortened, users := range shortenedUserSet {

		// Check if the shortened URL should be deleted.
		if len(users) == 0 {
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

func (s *shortenedIndex) lock(f func()) {
	s.mux.Lock()
	defer s.mux.Unlock()
	f()
}

func (s *shortenedIndex) rlock(f func()) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	f()
}

func (s *shortenedIndex) read(shortenedURLs []string) (shortenedUserSet map[string]userSet) {

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
			users := s.indexMap[shortened]

			// Add the users to the return map.
			shortenedUserSet[shortened] = users
		}
	}

	return shortenedUserSet
}
