package auth

import (
	"sync"
)

type userSet map[string]struct{}

// TODO
type shortenedIndex struct { // AKA authorization data structure 2.
	lookup map[string]userSet
	mux    sync.RWMutex
}

func newShortenedLookup() (shortLookup *shortenedIndex) {
	return &shortenedIndex{
		lookup: make(map[string]userSet),
	}
}

func (s *shortenedIndex) add(shortenedUserSet map[string]userSet) {

	// Iterate through the given shortened URL user sets.
	for shortened, users := range shortenedUserSet {

		// Confirm the shortened URL exists.
		_, ok := s.lookup[shortened]
		if !ok {
			s.lookup[shortened] = userSet{}
		}

		// Iterate through the given users and add them to the set.
		for user := range users {
			s.lookup[shortened][user] = struct{}{}
		}
	}
}

func (s *shortenedIndex) delete(shortenedUserSet map[string]userSet) {

	// Iterate through the given shortened URLs.
	for shortened, users := range shortenedUserSet {

		// Check if the shortened URL should be deleted.
		if len(users) == 0 {
			delete(s.lookup, shortened)
			continue
		}

		// Confirm the shortened URL is present in data structure 2.
		_, ok := s.lookup[shortened]
		if !ok {
			continue
		}

		// Delete the users.
		for user := range users {
			delete(s.lookup[shortened], user)
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

	// TODO Handle an empty case?

	// Iterate through the given shortened URLs.
	for _, shortened := range shortenedURLs {

		// Get the users associated with this shortened URL.
		users := s.lookup[shortened]

		// Add the users to the return map.
		shortenedUserSet[shortened] = users
	}

	return shortenedUserSet
}
