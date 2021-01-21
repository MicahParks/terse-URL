package storage

import (
	"context"

	"github.com/MicahParks/terseurl/models"
)

// InitializeSummaries TODO
func InitializeSummaries(ctx context.Context, terseStore TerseStore, visitsStore VisitsStore) (summaries map[string]models.TerseSummary, err error) {

	summaries = make(map[string]models.TerseSummary)

	//
	var counts map[string]uint
	if counts, err = visitsStore.ExportCounts(ctx); err != nil {
		return nil, err
	}

	// TODO Use ctxCreator function.

	var terse map[string]*models.Terse
	if terse, err = terseStore.ExportTerse(ctx); err != nil {
		return nil, err
	}

	for shortened, count := range counts {

		// TODO Check if key in terse not in counts?
		t, ok := terse[shortened]
		if !ok {
			//TODO
		}

		summaries[shortened] = models.TerseSummary{
			OriginalURL:  *t.OriginalURL, // TODO Pointers
			RedirectType: t.MediaPreview.RedirectType,
			ShortenedURL: *t.ShortenedURL, // TODO Pointers.
			VisitCount:   int64(count),    // TODO Int conversion.
		}
	}

	return summaries, nil
}
