package storage

import (
	"context"
	"sync"

	"github.com/MicahParks/terseurl/models"
)

// MemTerse is a TerseStore implementation that stores all data in a Go map in memory.
type MemTerse struct {
	mux   sync.RWMutex
	terse map[string]*models.Terse
}

// NewMemTerse creates a new MemTerse given the required assets.
func NewMemTerse() (terseStore TerseStore) {
	return &MemTerse{
		terse: make(map[string]*models.Terse),
	}
}

// Close closes the connection to the underlying storage.
func (m *MemTerse) Close(_ context.Context) (err error) {

	// Lock the Terse data for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Delete all the Terse data.
	m.deleteAll()

	return nil
}

// Delete deletes the Terse data for the given shortened URLs. If shortenedURLs is nil or empty, all shortened URL
// Terse data are deleted. There should be no error if a shortened URL is not found.
func (m *MemTerse) Delete(_ context.Context, shortenedURLs []string) (err error) {

	// Lock the Terse data for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Check for the empty case.
	if len(shortenedURLs) == 0 {

		// Delete all Terse data.
		m.deleteAll()
	} else {

		// Iterate through the given shortened URLs.
		for _, shortened := range shortenedURLs {
			delete(m.terse, shortened)
		}
	}

	return nil
}

// Read returns a map of shortened URLs to Terse data. If shortenedURLs is nil or empty, all shortened URL Terse
// data are expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
func (m *MemTerse) Read(_ context.Context, shortenedURLs []string) (terseData map[string]*models.Terse, err error) {

	// Create the return map.
	terseData = make(map[string]*models.Terse, len(shortenedURLs))

	// Lock the Terse data for async safe usage.
	m.mux.RLock()
	defer m.mux.RUnlock()

	// Check for the empty case.
	if len(shortenedURLs) == 0 {

		// Use all Terse data.
		terseData = m.terse
	} else {

		// Iterate through the given shortened URLs.
		for _, shortened := range shortenedURLs {

			// Get the Terse data for the shortened URL.
			terse, ok := m.terse[shortened]
			if !ok {
				return nil, ErrShortenedNotFound
			}

			// Add the Terse data to the return map.
			terseData[shortened] = terse
		}
	}

	return terseData, nil
}

// Summary summarizes the Terse data for the given shortened URLs. If shortenedURLs is nil or empty, then all
// shortened URL Summary data are expected.
func (m *MemTerse) Summary(_ context.Context, shortenedURLs []string) (summaries map[string]*models.TerseSummary, err error) {

	// Create the return map.
	summaries = make(map[string]*models.TerseSummary, len(shortenedURLs))

	// Lock the Terse data for async safe use.
	m.mux.RLock()
	defer m.mux.RUnlock()

	// Check for the empty case.
	if len(shortenedURLs) == 0 {

		// Gather the Summary data for all shortened URLs.
		for shortened, terse := range m.terse {
			summaries[shortened] = summarizeTerse(*terse)
		}
	} else {

		// Iterate through the given shortened URLs.
		for _, shortened := range shortenedURLs {

			// Get the Terse data for the shortened URL.
			terse, ok := m.terse[shortened]
			if !ok {
				return nil, ErrShortenedNotFound
			}

			// Add the Terse data to the return map.
			summaries[shortened] = summarizeTerse(*terse)
		}
	}

	return summaries, nil
}

// Write writes the given Terse data according to the given operation. The error must be storage.ErrShortenedExists
// if an Insert operation cannot be performed due to the Terse data already existing. The error must be
// storage.ErrShortenedNotFound if an Update operation cannot be performed due to the Terse data not existing.
func (m *MemTerse) Write(_ context.Context, terseData map[string]*models.Terse, operation WriteOperation) (err error) {

	// Lock the Terse data for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Iterate through the given Terse data.
	for shortened, terse := range terseData {

		// Check to see if the shortened URL already exists.
		if operation == Insert || operation == Update {
			_, ok := m.terse[shortened]
			if ok && operation == Insert {
				return ErrShortenedExists
			}
			if !ok && operation == Update {
				return ErrShortenedNotFound
			}
		}

		// Assign the shortened URL the given Terse data.
		m.terse[shortened] = terse
	}

	return nil
}

// deleteAll deletes all of the Terse data. It does not lock, so a lock must be used for async safe usage.
func (m *MemTerse) deleteAll() {

	// Reassign the Terse data so it's taken by the garbage collector.
	m.terse = make(map[string]*models.Terse)
}
