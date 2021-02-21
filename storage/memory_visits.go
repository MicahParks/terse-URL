package storage

import (
	"context"
	"sync"

	"github.com/MicahParks/terseurl/models"
)

// MemVisits is a VisitsStore implementation that stores all data in a Go map in memory.
type MemVisits struct {
	mux    sync.RWMutex
	visits map[string][]*models.Visit
}

// NewMemVisits creates a new MemVisits.
func NewMemVisits() (visitsStore VisitsStore) {
	return &MemVisits{
		visits: make(map[string][]*models.Visit),
	}
}

// Close closes the connection to the underlying storage.
func (m *MemVisits) Close(_ context.Context) (err error) {

	// Lock the Summary data for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Delete all the Summary data.
	m.deleteAll()

	return nil
}

// Delete deletes Visits data for the given shortened URLs. No error should be given if a shortened URL is not
// found.
func (m *MemVisits) Delete(_ context.Context, shortenedURLs []string) (err error) {

	// Lock the Summary data for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Check if all Summary data should be deleted.
	if shortenedURLs == nil {
		m.deleteAll()
		return nil
	}

	// Iterate through the given shortened URLs.
	for _, shortened := range shortenedURLs {
		delete(m.visits, shortened)
	}

	return nil
}

// Insert inserts the given Visits data. The visits do not need to be unique, so the Visits data should be appended
// to the data structure in storage.
func (m *MemVisits) Insert(_ context.Context, visitsData map[string][]*models.Visit) (err error) {

	// Lock the Summary data for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

}

// Read exports the Visits data for the given shortened URLs.
func (m *MemVisits) Read(_ context.Context, shortenedURLs []string) (visitsData map[string][]*models.Visit, err error) {

}

// Summary summarizes the Visits data for the given shortened URLs. This is used in building the SummaryStore.
func (m *MemVisits) Summary(_ context.Context, shortenedURLs []string) (visitsSummary map[string]*models.VisitsSummary, err error) {

}

// deleteAll deletes all of the Visits data. It does not lock, so a lock must be used for async safe usage.
func (m *MemVisits) deleteAll() {

	// Reassign the Visits data so it's taken by the garbage collector.
	m.visits = make(map[string][]*models.Visit)
}
