package storage

import (
	"context"

	"github.com/MicahParks/terseurl/models"
)

// SummaryStore is the Terse summary data storage interface. It allows for Terse summary storage operations without
// needing to know how the Terse summary data is stored.
type SummaryStore interface {

	// Close closes the connection to the underlying storage.
	Close(ctx context.Context) (err error)

	// Delete deletes the summary information for the given shortened URLs. If shortenedURLs is nil, all Summary data
	// will be deleted.
	Delete(ctx context.Context, shortenedURLs []string) (err error)

	// IncrementVisitCount increments the visit count for the given shortened URL. It is called in separate goroutine.
	// The error must be storage.ErrShortenedNotFound if the shortened URL is not found.
	IncrementVisitCount(ctx context.Context, shortened string) (err error)

	// Summary provides the summary information for the given shortened URLs. If shortenedURLs is nil, all summaries
	// will be returned. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
	Summary(ctx context.Context, shortenedURLs []string) (summaries map[string]models.Summary, err error)

	// Upsert upserts the summary information for the given shortened URL.
	Upsert(ctx context.Context, summaries map[string]models.Summary) (err error)
}

// TerseStore is the Terse storage interface. It allows for Terse storage operations without needing to know how
// the Terse data is stored.
type TerseStore interface {

	// Close closes the connection to the underlying storage.
	Close(ctx context.Context) (err error)

	// Delete deletes the Terse data for the given shortened URLs. If shortenedURLs is nil, all shortened URLs' Terse
	// data is deleted. There should be no error if a shortened URL is not found.
	Delete(ctx context.Context, shortenedURLs []string) (err error)

	// Read returns a map of shortened URLs to Terse data. If shortenedURLs is nil, all shortened URLs' Terse data is
	// expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
	Read(ctx context.Context, shortenedURLs []string) (terseData map[string]models.Terse, err error)

	// Summary summarizes the Terse data for the given shortened URLs. This is used in building the SummaryStore.
	Summary(ctx context.Context, shortenedURLs []string) (summaries map[string]models.TerseSummary, err error)

	// Write writes the given Terse data according to the given operation.
	Write(ctx context.Context, terseData map[string]models.Terse, operation WriteOperation) (err error)
}

// VisitsStore is the Visits storage interface. It allows for Visits storage operations without needing to know how the
// Visits data is stored.
type VisitsStore interface {

	// Close closes the connection to the underlying storage.
	Close(ctx context.Context) (err error)

	// Delete deletes Visits data for the given shortened URLs. No error should be given if a shortened URL is not
	// found.
	Delete(ctx context.Context, shortenedURLs []string) (err error)

	// Read exports the Visits data for the given shortened URLs.
	Read(ctx context.Context, shortenedURLs []string) (visitsData map[string][]*models.Visit, err error)

	// Insert inserts the given Visits data. The visits do not need to be unique, so the Visits data should be appended
	// to the data structure in storage.
	Insert(ctx context.Context, visitsData map[string][]*models.Visit) (err error)

	// Summary summarizes the Visits data for the given shortened URLs. This is used in building the SummaryStore.
	Summary(ctx context.Context, shortenedURLs []string) (visitsSummary map[string]uint, err error)
}
