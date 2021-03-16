package storage

import (
	"context"
	"sync"

	"github.com/MicahParks/terseurl/models"
)

// MemVisits is a VisitsStore implementation that stores all data in a Go map in memory.
type MemVisits struct {
	mux    sync.RWMutex
	visits map[string][]models.Visit
}

// NewMemVisits creates a new MemVisits.
func NewMemVisits() (visitsStore VisitsStore) {
	return &MemVisits{
		visits: make(map[string][]models.Visit),
	}
}

// Close closes the connection to the underlying storage.
func (m *MemVisits) Close(_ context.Context) (err error) {

	// Lock the Visits data for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Delete all the Visits data.
	m.deleteAll()

	return nil
}

// Delete deletes Visits data for the given shortened URLs. If shortenedURLs is nil or empty, then all Visits data
// are deleted. No error should be given if a shortened URL is not found.
func (m *MemVisits) Delete(_ context.Context, shortenedURLs []string) (err error) {

	// Lock the Visits data for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Check for the empty case.
	if len(shortenedURLs) == 0 {

		// Delete all Visits data.
		m.deleteAll()
	} else {

		// Iterate through the given shortened URLs.
		for _, shortened := range shortenedURLs {
			delete(m.visits, shortened)
		}
	}

	return nil
}

// Insert inserts the given Visits data. The visits do not need to be unique, so the Visits data should be appended
// to the data structure in storage.
func (m *MemVisits) Insert(_ context.Context, visitsData map[string][]models.Visit) (err error) {

	// Lock the Visits data for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Iterate through the given Visits data and append it to the existing.
	for shortened, visits := range visitsData {
		visitsData[shortened] = append(visitsData[shortened], visits...)
	}

	return nil
}

// Read exports the Visits data for the given shortened URLs. If shortenedURLs is nil or empty, then all shortened
// URL Visits data are expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
func (m *MemVisits) Read(_ context.Context, shortenedURLs []string) (visitsData map[string][]models.Visit, err error) {

	// Create the return map.
	visitsData = make(map[string][]models.Visit, len(shortenedURLs))

	// Lock the Visits data for async safe use.
	m.mux.RLock()
	defer m.mux.RUnlock()

	// Check for the empty case.
	if len(shortenedURLs) == 0 {

		// Use all Visits data.
		visitsData = m.visits
	} else {

		// Iterate through the given shortened URLs.
		for _, shortened := range shortenedURLs {

			// Get the Visits data for the shortened URL.
			visits, ok := m.visits[shortened]
			if !ok {
				return nil, ErrShortenedNotFound
			}

			// Add the Visits data to the return map.
			visitsData[shortened] = visits
		}
	}

	return visitsData, nil
}

// Summary summarizes the Visits data for the given shortened URLs. If shortenedURLs is nil or empty, then all
// shortened URL Summary data are expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not
// found.
func (m *MemVisits) Summary(_ context.Context, shortenedURLs []string) (summaries map[string]*models.VisitsSummary, err error) {

	// Create the return map.
	summaries = make(map[string]*models.VisitsSummary, len(shortenedURLs))

	// Lock the Visits data for async safe use.
	m.mux.RLock()
	defer m.mux.RUnlock()

	// Check for the empty case.
	if len(shortenedURLs) == 0 {

		// Gather the Summary data for all shortened URLs.
		for shortened, visits := range m.visits {
			summaries[shortened] = &models.VisitsSummary{
				VisitCount: uint64(len(visits)),
			}
		}
	} else {

		// Iterate through the given shortened URLs.
		for _, shortened := range shortenedURLs {

			// Get the Visits data for the shortened URL.
			visits, ok := m.visits[shortened]
			if !ok {
				return nil, ErrShortenedNotFound
			}

			// Add the Visits data to the return map.
			summaries[shortened] = &models.VisitsSummary{
				VisitCount: uint64(len(visits)),
			}
		}
	}

	return summaries, nil
}

// deleteAll deletes all of the Visits data. It does not lock, so a lock must be used for async safe usage.
func (m *MemVisits) deleteAll() {

	// Reassign the Visits data so it's taken by the garbage collector.
	m.visits = make(map[string][]models.Visit)
}
