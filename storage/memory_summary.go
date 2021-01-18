package storage

import (
	"context"

	"github.com/MicahParks/terseurl/models"
)

// TODO
type MemSummary struct {
	summaries map[string]*models.TerseSummary
}

// TODO
func NewMemSummary(visitsStore VisitsStore) SummaryStore {
	// TODO Count all visits in VisitsStore. Maybe make a method for it.
	return &MemSummary{
		summaries: make(map[string]*models.TerseSummary),
	}
}

// TODO
func (m *MemSummary) IncrementVisitCount(_ context.Context, shortened string) (err error) {

	// Get the summary.
	summary, ok := m.summaries[shortened]
	if !ok {
		// TODO Return error.
	}

	// Increment the visit counter.
	summary.VisitCount++

	return nil
}

// TODO
func (m *MemSummary) Summarize(_ context.Context, shortenedURLs []string) (summaries map[string]*models.TerseSummary, err error) {

	// Make a map of summaries.
	summaries = make(map[string]*models.TerseSummary)

	// Iterate through the shortened URLs.
	for _, shortened := range shortenedURLs {
		summary, ok := m.summaries[shortened]
		if !ok {
			continue // TODO Return error.
		}
		summaries[shortened] = summary
	}

	return summaries, nil
}

// TODO
func (m *MemSummary) Upsert(ctx context.Context, shortened string, summary *models.TerseSummary) (err error) {
	// TODO
	return
}
