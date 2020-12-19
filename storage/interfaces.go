package storage

import (
	"context"

	"github.com/MicahParks/terse-URL/models"
)

// TerseStore is the Terse storage interface. It allows for Terse storage operations without needing to know how
// the Terse data is stored.
type TerseStore interface {

	// Close closes the connection to the underlying storage. The ctxerrgroup should be killed. This may or may not
	// close the connection to the VisitsStore, depending on the configuration.
	Close(ctx context.Context) (err error)

	// Import imports the given export's data. If del is not nil, data will be deleted accordingly. If del is nil, data
	// may be overwritten, but unaffected data will be untouched. If the VisitsStore is not nil, then the same method
	// will be called for the associated VisitsStore.
	Import(ctx context.Context, delete *models.Delete, export map[string]*models.Export) (err error)

	// InsertTerse adds a Terse to the TerseStore. The shortened URL will be active after this. The error must be
	// storage.ErrShortenedExists if the shortened URL is already present.
	InsertTerse(ctx context.Context, terse *models.Terse) (err error)

	// Delete deletes data according to the del argument. If the VisitsStore is not nil, then the same method will be
	// called for the associated VisitsStore.
	Delete(ctx context.Context, delete models.Delete) (err error)

	// DeleteTerse deletes data according to the del argument for the given shortened URL. No error should be given if
	// the shortened URL is not found. If the VisitsStore is not nil, then the same method will be called for the
	// associated VisitsStore.
	DeleteTerse(ctx context.Context, del models.Delete, shortened string) (err error)

	// Export returns a export of Terse and Visit data for a given shortened URL. The error must be
	// storage.ErrShortenedNotFound if the shortened URL is not found.
	Export(ctx context.Context, shortened string) (export models.Export, err error)

	// ExportAll returns a map of shortened URLs to export data.
	ExportAll(ctx context.Context) (export map[string]models.Export, err error)

	// ReadTerse retrieves all non-Visit Terse data give its shortened URL. A nil visit may be passed in and the visit
	// should not be recorded. The error must be storage.ErrShortenedNotFound if the shortened URL is not found.
	ReadTerse(ctx context.Context, shortened string, visit *models.Visit) (terse *models.Terse, err error)

	// UpdateTerse assumes the Terse already exists. It will override all of its values. The error must be
	// storage.ErrShortenedNotFound if the shortened URL is not found.
	UpdateTerse(ctx context.Context, terse *models.Terse) (err error)

	// UpsertTerse will upsert the Terse into the backend storage.
	UpsertTerse(ctx context.Context, terse *models.Terse) (err error)

	// VisitsStore returns the underlying VisitsStore, which hold the backend storage for tracking visits to shortened
	// URLs.
	VisitsStore() VisitsStore
}

// VisitsStore is the Visits storage interface. It allows for Visits storage operations without needing to know how the
// Visits data is stored.
type VisitsStore interface {

	// CreateVisit adds the visit to the visits store.
	AddVisit(ctx context.Context, shortened string, visit *models.Visit) (err error)

	// Close closes the connection to the underlying storage.
	Close(ctx context.Context) (err error)

	// Delete deletes data according to the del argument.
	Delete(ctx context.Context, del models.Delete) (err error)

	// DeleteVisits deletes data according to the del argument for the shortened URL. No error should be given if the
	// shortened URL is not found.
	DeleteVisits(ctx context.Context, del models.Delete, shortened string) (err error)

	// All exports all visits data.
	Export(ctx context.Context) (allVisits map[string][]*models.Visit, err error)

	// ExportShortened gets all visits to the shortened URL. The error must be storage.ErrShortenedNotFound if the
	// shortened URL is not found.
	ExportShortened(ctx context.Context, shortened string) (visits []*models.Visit, err error)

	// Import imports the given export's data. If del is not nil, data will be deleted accordingly. If del is nil, data
	// may be overwritten, but unaffected data will be untouched.
	Import(ctx context.Context, del *models.Delete, export map[string]*models.Export) (err error)
}
