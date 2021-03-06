package auth

import (
	"sync"
)

// TODO
type shortenedLookup struct { // AKA authorization data structure 2.
	lookup map[string]map[string]struct{}
	mux    sync.RWMutex
}

func newShortenedLookup() (shortLookup *shortenedLookup) {
	return &shortenedLookup{
		lookup: make(map[string]map[string]struct{}),
	}
}

func (s *shortenedLookup) add(shortened string, users []string) {

	// Confirm the shortened URL exists.
	_, ok := s.lookup[shortened]
	if !ok {
		s.lookup[shortened] = make(map[string]struct{})
	}

	// Iterate through the given users and add them to the set.
	for _, user := range users {
		s.lookup[shortened][user] = struct{}{}
	}
}

func (s *shortenedLookup) delete(shortenedUsers map[string][]string) {

	// Iterate through the given shortened URLs.
	for shortened, users := range shortenedUsers {

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
		for _, user := range users {
			delete(s.lookup[shortened], user)
		}
	}
}

func (s *shortenedLookup) lock(f func()) {
	s.mux.Lock()
	defer s.mux.Unlock()
	f()
}

func (s *shortenedLookup) rlock(f func()) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	f()
}

func (s *shortenedLookup) read(shortenedURLs []string) (shortenedUsers map[string]map[string]struct{}) {

	// Create the return map.
	shortenedUsers = make(map[string]map[string]struct{})

	// TODO Handle an empty case?

	// Iterate through the given shortened URLs.
	for _, shortened := range shortenedURLs {

		// Get the users associated with this shortened URL.
		users := s.lookup[shortened]

		// Add the users to the return map.
		shortenedUsers[shortened] = users
	}

	return shortenedUsers
}
