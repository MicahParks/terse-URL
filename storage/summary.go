package storage

import (
	"context"

	"github.com/MicahParks/terseurl/models"
)

// InitializeSummaries is a helper function to create the required data for a SummaryStore on startup. It will iterate
// through all assets in the TerseStore and VisitsStore and map their shortened URLs to Terse summary data.
func InitializeSummaries(ctx context.Context, terseStore TerseStore, visitsStore VisitsStore) (summaries map[string]models.TerseSummary, err error) {

	// Create the return map.
	summaries = make(map[string]models.TerseSummary)

	// Create the map of counts.
	counts := make(map[string]uint)
	if visitsStore != nil {
		if counts, err = visitsStore.ExportCounts(ctx); err != nil {
			return nil, err
		}
	}

	// Create the map of Terse data.
	var terse map[string]*models.Terse
	if terse, err = terseStore.ExportTerse(ctx); err != nil {
		return nil, err
	}

	// Iterate through the shortened URLs in the Terse data.
	for shortened, t := range terse {

		// Get the count for the shortened URL.
		count := counts[shortened]

		// Combine the required Terse data with the count of Visits to create Terse summary data.
		summaries[shortened] = models.TerseSummary{
			OriginalURL:  t.OriginalURL,
			RedirectType: t.RedirectType,
			ShortenedURL: t.ShortenedURL,
			VisitCount:   int64(count), // TODO Conversion may cause data loss.
		}
	}

	return summaries, nil
}
