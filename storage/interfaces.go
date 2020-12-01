package storage

import (
	"context"
	"time"

	"github.com/MicahParks/terse-URL/models"
)

// TerseStore is the shortened link storage interface. It allows for link storage operations without needing to know how
// or where the shortened link data is stored.
type TerseStore interface { // TODO Add a method that removes all redirections based on original?

	// Close closes the connection to the underlying storage. This does not close the connection to the VisitsStore.
	Close(ctx context.Context) (err error)

	// RescheduleDeletions deletes expired links in the underlying storage and schedules deletions for the rest. It
	// should be called on startup for underlying storage that is persistent.
	ScheduleDeletions(ctx context.Context) (err error)

	// DeleteLink deletes the given shortened URL.
	DeleteTerse(ctx context.Context, shortened string) (err error)

	// GetLink retrieves an original URL given it's shortened URL. The visit will be stored in the VisitsStore, unless it
	// is nil. visitCancel and visitCtx are the context.CancelFunc and context.Context for the VisitsStore interactions.
	GetTerse(ctx context.Context, shortened string, visit *models.Visit, visitCancel context.CancelFunc, visitCtx context.Context) (original string, err error)

	// PutLink inserts the original and shortened URL key value pair to the underlying storage. The shortened link is
	// considered active after this. If deleteAt is nil, the Terse will not be scheduled for deletion.
	UpsertTerse(ctx context.Context, deleteAt *time.Time, original, shortened string) (err error)

	// VisitsStore returns the underlying VisitsStore, which hold the backend storage for tracking visits to shortened
	// URLs.
	VisitsStore() VisitsStore
}

// VisitsStore is the shortened link storage interface. It allows for storage operations relating to tracking shortened
// URL visits without needing to know where the visits are stored.
type VisitsStore interface {

	// Close closes the connection to the underlying storage.
	Close(ctx context.Context) (err error)

	// DeleteVisits deletes all visits to the shortened URL.
	DeleteVisits(ctx context.Context, shortened string) (err error)

	// GetVisits gets all visits to the shortened URL.
	GetVisits(ctx context.Context, shortened string) (visits []*models.Visit, err error)

	// PutVisit stores the visit to the shortened URL.
	AddVisit(ctx context.Context, shortened string, visit *models.Visit) (err error)
}
