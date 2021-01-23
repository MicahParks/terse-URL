package storage

import (
	"context"

	"github.com/MicahParks/terseurl/models"
)

// InitializeSummaries TODO
func InitializeSummaries(ctx context.Context, terseStore TerseStore, visitsStore VisitsStore) (summaries map[string]models.TerseSummary, err error) {

	// Create the return map.
	summaries = make(map[string]models.TerseSummary)

	//
	counts := make(map[string]uint)
	if visitsStore != nil {
		if counts, err = visitsStore.ExportCounts(ctx); err != nil {
			return nil, err
		}
	}

	//
	var terse map[string]*models.Terse
	if terse, err = terseStore.ExportTerse(ctx); err != nil {
		return nil, err
	}

	//
	for shortened, t := range terse {

		// TODO
		count, ok := counts[shortened]
		if !ok {
			//TODO
		}

		// TODO
		summaries[shortened] = models.TerseSummary{
			OriginalURL:  t.OriginalURL, // TODO Pointers
			RedirectType: t.RedirectType,
			ShortenedURL: t.ShortenedURL, // TODO Pointers.
			VisitCount:   int64(count),   // TODO Int conversion.
		}
	}

	return summaries, nil
}
