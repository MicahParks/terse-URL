package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/MicahParks/terseurl/models"
)

// MemSummary is a SummaryStore implementation that stores all data in a Go map in memory.
type MemSummary struct {
	summaries map[string]*models.TerseSummary
	mux       sync.RWMutex
}

// NewMemSummary creates a new MemSummary.
func NewMemSummary(visitsStore VisitsStore) (summaryStore SummaryStore) {
	// TODO Count all visits in VisitsStore. Maybe make a method for it.
	return &MemSummary{
		summaries: make(map[string]*models.TerseSummary),
	}
}

// IncrementVisitCount increments the visit count for the given shortened URL. It is called in separate goroutine. This
// implementation has no network activity and ignores the given context.
func (m *MemSummary) IncrementVisitCount(_ context.Context, shortened string) (err error) {

	// Lock the Summary map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Get the summary.
	summary, ok := m.summaries[shortened]
	if !ok {
		return fmt.Errorf("%w: %s", ErrShortenedNotFound, shortened)
	}

	// Increment the visit counter.
	summary.VisitCount++

	return nil
}

// Summarize provides the summary information for the given shortened URLs. This implementation has no network activity
// and ignores the given context.
func (m *MemSummary) Summarize(_ context.Context, shortenedURLs []string) (summaries map[string]models.TerseSummary, err error) {

	// Lock the Summary map for async safe use.
	m.mux.RLock()
	defer m.mux.RUnlock()

	// Make a map of summaries.
	summaries = make(map[string]models.TerseSummary)

	// Iterate through the shortened URLs. Add their summaries to the map.
	for _, shortened := range shortenedURLs {
		summary, ok := m.summaries[shortened]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrShortenedNotFound, shortened)
		}
		summaries[shortened] = *summary // TODO Should underlying map have pointers?
		// TODO Should there be pointers for a lot of these underlying data types?
	}

	return summaries, nil
}

// // Upsert upserts the summary information for the given shortened URL. This implementation has no network activity
// and ignores the given context.
func (m *MemSummary) Upsert(_ context.Context, summaries map[string]*models.TerseSummary) (err error) {

	// Lock the Summary map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// TODO

	return
}
