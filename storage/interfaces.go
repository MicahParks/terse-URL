package storage

import (
	"context"

	"github.com/MicahParks/terseurl/models"
)

// SummaryStore is the Terse summary data storage interface. It allows for Terse summary storage operations without
// needing to know how the Terse summary data is stored.
type SummaryStore interface {

	// IncrementVisitCount increments the visit count for the given shortened URL. It is called in separate goroutine.
	// The error must be storage.ErrShortenedNotFound if the shortened URL is not found.
	IncrementVisitCount(ctx context.Context, shortened string) (err error)

	// Summarize provides the summary information for the given shortened URLs. The error must be
	// storage.ErrShortenedNotFound if the shortened URL is not found.
	Summarize(ctx context.Context, shortenedURLs []string) (summaries map[string]models.TerseSummary, err error)

	// Upsert upserts the summary information for the given shortened URL.
	Upsert(ctx context.Context, summaries map[string]models.TerseSummary) (err error)
}

// TerseStore is the Terse storage interface. It allows for Terse storage operations without needing to know how
// the Terse data is stored.
type TerseStore interface {

	// Close closes the connection to the underlying storage. The ctxerrgroup should be killed. This may or may not
	// close the connection to the VisitsStore, depending on the configuration.
	Close(ctx context.Context) (err error)

	// CreateSummaryStore creates the SummaryStore based on the existing VisitsStore data.
	CreateSummaryStore() (summaries map[string]models.TerseSummary, err error)

	// Delete deletes data according to the del argument. If the VisitsStore is not nil, then the same method will be
	// called for the associated VisitsStore.
	Delete(ctx context.Context, del models.Delete) (err error)

	// DeleteOne deletes data according to the del argument for the given shortened URL. No error should be given if
	// the shortened URL is not found. If the VisitsStore is not nil, then the same method will be called for the
	// associated VisitsStore.
	DeleteOne(ctx context.Context, del models.Delete, shortened string) (err error)

	// Export returns a map of shortened URLs to export data.
	Export(ctx context.Context) (export map[string]models.Export, err error)

	// ExportOne returns a export of Terse and Visit data for a given shortened URL. The error must be
	// storage.ErrShortenedNotFound if the shortened URL is not found.
	ExportOne(ctx context.Context, shortened string) (export models.Export, err error)

	// Import imports the given export's data. If del is not nil, data will be deleted accordingly. If del is nil, data
	// may be overwritten, but unaffected data will be untouched. If the VisitsStore is not nil, then the same method
	// will be called for the associated VisitsStore.
	Import(ctx context.Context, del *models.Delete, export map[string]models.Export) (err error)

	// Insert adds a Terse to the TerseStore. The shortened URL will be active after this. The error must be
	// storage.ErrShortenedExists if the shortened URL is already present.
	Insert(ctx context.Context, terse *models.Terse) (err error)

	// Read retrieves all non-Visit Terse data give its shortened URL. A nil visit may be passed in and the visit should
	// not be recorded. The error must be storage.ErrShortenedNotFound if the shortened URL is not found.
	Read(ctx context.Context, shortened string, visit *models.Visit) (terse *models.Terse, err error)

	// SummaryStore returns the underlying SummaryStore, which holds the backend storage for tracking summary data for
	// Terse.
	SummaryStore() SummaryStore

	// Update assumes the Terse already exists. It will override all of its values. The error must be
	// storage.ErrShortenedNotFound if the shortened URL is not found.
	Update(ctx context.Context, terse *models.Terse) (err error)

	// Upsert will upsert the Terse into the backend storage.
	Upsert(ctx context.Context, terse *models.Terse) (err error)

	// VisitsStore returns the underlying VisitsStore, which holds the backend storage for tracking visits to shortened
	// URLs.
	VisitsStore() VisitsStore
}

// VisitsStore is the Visits storage interface. It allows for Visits storage operations without needing to know how the
// Visits data is stored.
type VisitsStore interface {

	// Add adds the visit to the visits store.
	Add(ctx context.Context, shortened string, visit *models.Visit) (err error)

	// Close closes the connection to the underlying storage.
	Close(ctx context.Context) (err error)

	// Delete deletes data according to the del argument.
	Delete(ctx context.Context, del models.Delete) (err error)

	// DeleteOne deletes data according to the del argument for the shortened URL. No error should be given if the
	// shortened URL is not found.
	DeleteOne(ctx context.Context, del models.Delete, shortenedURLs []string) (err error)

	// Export exports all exports all visits data.
	Export(ctx context.Context) (allVisits map[string][]*models.Visit, err error)

	// ExportCounts creates a map of shortened URLs to count of Visits.
	ExportCounts(ctx context.Context) (counts map[string]uint, err error)

	// ExportSome gets all visits to the shortened URL. The error must be storage.ErrShortenedNotFound if the shortened
	// URL is not found.
	ExportSome(ctx context.Context, shortenedURLs []string) (visits map[string][]*models.Visit, err error)

	// Import imports the given export's data. If del is not nil, data will be deleted accordingly. If del is nil, data
	// may be overwritten, but unaffected data will be untouched.
	Import(ctx context.Context, del *models.Delete, export map[string]models.Export) (err error)
}
