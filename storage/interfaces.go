package storage

import (
	"context"

	"github.com/MicahParks/terseurl/models"
)

// SummaryStore is the Summary data storage interface. It allows for Summary data storage operations without needing to
// know how the Summary data are stored. Summary data should not persist through a service restart.
type SummaryStore interface {

	// Close closes the connection to the underlying storage.
	Close(ctx context.Context) (err error)

	// Delete deletes the summary information for the given shortened URLs. If shortenedURLs is nil or empty, all
	// Summary data are deleted. No error should be returned if a shortened URL is not found.
	Delete(ctx context.Context, shortenedURLs []string) (err error)

	// IncrementVisitCount increments the visit count for the given shortened URL. It is called in separate goroutine.
	// The error must be storage.ErrShortenedNotFound if the shortened URL is not found.
	IncrementVisitCount(ctx context.Context, shortened string) (err error)

	// Read provides the summary information for the given shortened URLs. If shortenedURLs is nil or empty, all
	// summaries are returned. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
	Read(ctx context.Context, shortenedURLs []string) (summaries map[string]*models.Summary, err error)

	// Upsert upserts the summary information for the given shortened URL.
	Upsert(ctx context.Context, summaries map[string]*models.Summary) (err error)
}

// TerseStore is the Terse storage interface. It allows for Terse storage operations without needing to know how
// the Terse data are stored.
type TerseStore interface {

	// Close closes the connection to the underlying storage.
	Close(ctx context.Context) (err error)

	// Delete deletes the Terse data for the given shortened URLs. If shortenedURLs is nil or empty, all shortened URL
	// Terse data are deleted. There should be no error if a shortened URL is not found.
	Delete(ctx context.Context, shortenedURLs []string) (err error)

	// Read returns a map of shortened URLs to Terse data. If shortenedURLs is nil or empty, all shortened URL Terse
	// data are expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
	Read(ctx context.Context, shortenedURLs []string) (terseData map[string]*models.Terse, err error)

	// Summary summarizes the Terse data for the given shortened URLs. If shortenedURLs is nil or empty, then all
	// shortened URL Summary data are expected.
	Summary(ctx context.Context, shortenedURLs []string) (summaries map[string]*models.TerseSummary, err error)

	// Write writes the given Terse data according to the given operation. The error must be storage.ErrShortenedExists
	// if an Insert operation cannot be performed due to the Terse data already existing. The error must be
	// storage.ErrShortenedNotFound if an Update operation cannot be performed due to the Terse data not existing.
	Write(ctx context.Context, terseData map[string]*models.Terse, operation WriteOperation) (err error)
}

// VisitsStore is the Visits storage interface. It allows for Visits storage operations without needing to know how the
// Visits data are stored.
type VisitsStore interface {

	// Close closes the connection to the underlying storage.
	Close(ctx context.Context) (err error)

	// Delete deletes Visits data for the given shortened URLs. If shortenedURLs is nil or empty, then all Visits data
	// are deleted. No error should be given if a shortened URL is not found.
	Delete(ctx context.Context, shortenedURLs []string) (err error)

	// Insert inserts the given Visits data. The visits do not need to be unique, so the Visits data should be appended
	// to the data structure in storage.
	Insert(ctx context.Context, visitsData map[string][]models.Visit) (err error)

	// Read exports the Visits data for the given shortened URLs. If shortenedURLs is nil or empty, then all shortened
	// URL Visits data are expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
	Read(ctx context.Context, shortenedURLs []string) (visitsData map[string][]models.Visit, err error)

	// Summary summarizes the Visits data for the given shortened URLs. If shortenedURLs is nil or empty, then all
	// shortened URL Summary data are expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not
	// found.
	Summary(ctx context.Context, shortenedURLs []string) (summaries map[string]*models.VisitsSummary, err error)
}
