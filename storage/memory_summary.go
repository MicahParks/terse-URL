package storage

import (
	"context"
	"sync"

	"github.com/MicahParks/terseurl/models"
)

// TODO Memory implementations for all stores both consume and produce pointers. Should copies of the data be consumed
// TODO and produced so the user doesn't have any confusion? (Slices as well.)

// MemSummary is a SummaryStore implementation that stores all data in a Go map in memory.
type MemSummary struct {
	summaries map[string]*models.Summary
	mux       sync.RWMutex
}

// NewMemSummary creates a new MemSummary.
func NewMemSummary() (summaryStore SummaryStore) {
	return &MemSummary{
		summaries: make(map[string]*models.Summary),
	}
}

// Close closes the connection to the underlying storage.
func (m *MemSummary) Close(_ context.Context) (err error) {

	// Lock the Summary data for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Delete all the Summary data.
	m.deleteAll()

	return nil
}

// Delete deletes the summary information for the given shortened URLs. If shortenedURLs is nil, all Summary data are
// deleted. No error should be returned if a shortened URL is not found.
func (m *MemSummary) Delete(_ context.Context, shortenedURLs []string) (err error) {

	// Lock the Summary data for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Check for the nil case.
	if shortenedURLs == nil {

		// Delete all summary data.
		m.deleteAll()
	} else {

		// Iterate through the given shortened URLs.
		for _, shortened := range shortenedURLs {
			delete(m.summaries, shortened)
		}
	}

	return nil
}

// IncrementVisitCount increments the visit count for the given shortened URL. It is called in separate goroutine.
// The error must be storage.ErrShortenedNotFound if the shortened URL is not found.
func (m *MemSummary) IncrementVisitCount(_ context.Context, shortened string) (err error) {

	// Lock the Summary data for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Get the current summary data.
	summary, ok := m.summaries[shortened]
	if !ok {
		return ErrShortenedNotFound
	}

	// Increment the visits count.
	summary.Visits.VisitCount++

	// Reassign the summary data.
	m.summaries[shortened] = summary

	return nil
}

// Read provides the summary information for the given shortened URLs. If shortenedURLs is nil, all summaries are
// returned. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
func (m *MemSummary) Read(_ context.Context, shortenedURLs []string) (summaries map[string]*models.Summary, err error) {

	// Create the return map.
	summaries = make(map[string]*models.Summary, len(shortenedURLs))

	// Lock the Summary data for async safe use.
	m.mux.RLock()
	defer m.mux.RUnlock()

	// Check to see if all Summary data was requested, if so, copy all Summary data.
	if shortenedURLs == nil {

		// Use all Summary data.
		summaries = m.summaries
	} else {

		// Iterate through the given shortened URLs. Copy the requested ones.
		for _, shortened := range shortenedURLs {
			summary, ok := m.summaries[shortened]
			if !ok {
				return nil, ErrShortenedNotFound
			}
			summaries[shortened] = summary
		}
	}

	return summaries, nil
}

// Upsert upserts the Summary data for the given shortened URL.
func (m *MemSummary) Upsert(_ context.Context, summaries map[string]*models.Summary) (err error) {

	// Lock the Summary data for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Iterate through the given summary data. Upsert the Summary data.
	for shortened, summary := range summaries {
		m.summaries[shortened] = summary
	}

	return nil
}

// deleteAll deletes all of the Summary data. It does not lock, so a lock must be used for async safe usage.
func (m *MemSummary) deleteAll() {

	// Reassign the Summary data so it's taken by the garbage collector.
	m.summaries = make(map[string]*models.Summary)
}
