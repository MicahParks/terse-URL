package storage

import (
	"context"

	"github.com/MicahParks/terse-URL/models"
)

// TerseStore is the Terse storage interface. It allows for Terse storage operations without needing to know how
// the Terse data is stored.
type TerseStore interface {

	// AddTerse adds a Terse to the TerseStore. The shortened URL will be active after this.
	AddTerse(ctx context.Context, terse *models.Terse)

	// Close closes the connection to the underlying storage. This does not close the connection to the VisitsStore.
	Close(ctx context.Context) (err error)

	// DeleteTerse deletes the given shortened URL.
	DeleteTerse(ctx context.Context, shortened string) (err error)

	// Dump returns a dump of all data for a given shortened URL.
	Dump(ctx context.Context, shortened string) (dump *models.Dump, err error)

	// DumpAll returns a map of shortened URLs to dump data.
	DumpAll(ctx context.Context) (dump map[string]*models.Dump, err error)

	// GetTerse retrieves all non-Visit Terse data give its shortened URL. The error must be
	// storage.ErrShortenedNotFound if the shortened URL is not found.
	// TODO Delete Terse if expired
	GetTerse(ctx context.Context, shortened string, visit *models.Visit) (original string, err error)

	// UpdateTerse assumes the Terse already exists. It will override all of its values. The error must be
	// storage.ErrShortenedNotFound if the shortened URL is not found.
	UpdateTerse(ctx context.Context, terse *models.Terse) (err error)

	// VisitsStore returns the underlying VisitsStore, which hold the backend storage for tracking visits to shortened
	// URLs.
	VisitsStore() VisitsStore
}

// VisitsStore is the Visits storage interface. It allows for Visits storage operations without needing to know how the
// Visits data is stored.
type VisitsStore interface {

	// AddVisit adds the visit to the visits store.
	AddVisit(ctx context.Context, shortened string, visit *models.Visit) (err error)

	// Close closes the connection to the underlying storage.
	Close(ctx context.Context) (err error)

	// DeleteVisits deletes all visits to the shortened URL.
	DeleteVisits(ctx context.Context, shortened string) (err error)

	// GetVisits gets all visits to the shortened URL.
	GetVisits(ctx context.Context, shortened string) (visits []*models.Visit, err error)
}
