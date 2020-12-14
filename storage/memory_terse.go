package storage

import (
	"context"
	"sync"

	"github.com/MicahParks/ctxerrgroup"

	"github.com/MicahParks/terse-URL/models"
)

// MemTerse is a TerseStore implementation that stores all data in a Go map in memory.
type MemTerse struct {
	createCtx   ctxCreator
	errChan     chan<- error
	group       *ctxerrgroup.Group
	mux         sync.RWMutex
	terse       map[string]*models.Terse
	visitsStore VisitsStore
}

// NewMemTerse creates a new MemTerse given the required assets.
func NewMemTerse(createCtx ctxCreator, errChan chan<- error, group *ctxerrgroup.Group, visitsStore VisitsStore) (terseStore TerseStore) {
	return &MemTerse{
		createCtx:   createCtx,
		errChan:     errChan,
		group:       group,
		terse:       make(map[string]*models.Terse),
		visitsStore: visitsStore,
	}
}

// Close lets the garbage collector take care of the old Terse data and kills the ctxerrgroup.
func (m *MemTerse) Close(_ context.Context) (err error) {

	// Kill the worker pool.
	m.group.Kill()

	// Let the garbage collector take care of the old Terse data.
	m.terse = make(map[string]*models.Terse)

	return nil
}

// DeleteTerse deletes the Terse data fro the given shortened URL. This implementation has no network activity and
// ignores the given context.
func (m *MemTerse) DeleteTerse(_ context.Context, shortened string) (err error) {

	// Lock the Terse map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Delete the Terse from the Terse map.
	delete(m.terse, shortened)

	return nil
}

// Export exports Terse and Visits data for the given shortened URL. This implementation has no network activity and
// partially ignores the given context.
func (m *MemTerse) Export(ctx context.Context, shortened string) (export models.Export, err error) {

	// Lock the Terse map for async safe use.
	m.mux.RLock()
	defer m.mux.RUnlock()

	// Get the Terse.
	terse, ok := m.terse[shortened]
	if !ok {
		return models.Export{}, ErrShortenedNotFound
	}

	// Get the visits for the Terse.
	visits := make([]*models.Visit, 0)
	if m.visitsStore != nil {
		if visits, err = m.visitsStore.ReadVisits(ctx, shortened); err != nil {
			return models.Export{}, err
		}
	}

	return models.Export{
		Terse:  terse,
		Visits: visits,
	}, nil
}

// ExportAll exports all Terse and Visits data. This implementation has no network activity and partially ignores the
// given context.
func (m *MemTerse) ExportAll(ctx context.Context) (export map[string]models.Export, err error) {

	// Lock the Terse map for async safe use.
	m.mux.RLock()
	defer m.mux.RUnlock()

	// Create the export map.
	export = make(map[string]models.Export)

	// Iterate through the shortened URLs and add them to the map.
	for shortened, terse := range m.terse {

		// Get the visits for the Terse.
		visits := make([]*models.Visit, 0)
		if m.visitsStore != nil {
			if visits, err = m.visitsStore.ReadVisits(ctx, *terse.ShortenedURL); err != nil {
				return nil, err
			}
		}

		// Add the shortened URL and its visits to the data export.
		export[shortened] = models.Export{
			Terse:  terse,
			Visits: visits,
		}
	}

	return export, nil
}

// InsertTerse inserts the Terse into the TerseStore. It will fail if the Terse already exists. This implementation has
// no network activity and ignores the given context.
func (m *MemTerse) InsertTerse(_ context.Context, terse *models.Terse) (err error) {

	// Lock the Terse map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Check to see if the shortened URL already exists.
	if _, ok := m.terse[*terse.ShortenedURL]; ok {
		return ErrShortenedExists
	}

	// Add the shortened URL to the Terse map.
	m.terse[*terse.ShortenedURL] = terse

	return nil
}

// ReadTerse gets the Terse data for the given shortened URL. This implementation has no network activity and ignores
// the given context.
func (m *MemTerse) ReadTerse(_ context.Context, shortened string, visit *models.Visit) (terse *models.Terse, err error) {

	// Track the visit to this shortened URL. Do this in a separate goroutine so the response is faster.
	if visit != nil && m.visitsStore != nil {
		ctx, cancel := m.createCtx()
		go m.group.AddWorkItem(ctx, cancel, func(workCtx context.Context) (err error) {
			return m.visitsStore.AddVisit(workCtx, shortened, visit)
		})
	}

	// Lock the Terse map for async safe use.
	m.mux.RLock()
	defer m.mux.RUnlock()

	// Check to see if the shortened URL already exists.
	var ok bool
	terse, ok = m.terse[shortened]
	if !ok {
		return nil, ErrShortenedNotFound
	}

	return terse, nil
}

// UpdateTerse updates the Terse into the TerseStore. It will fail if the Terse does not exist. This implementation has
// no network activity and ignores the given context.
func (m *MemTerse) UpdateTerse(_ context.Context, terse *models.Terse) (err error) {

	// Lock the Terse map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Check to see if the shortened URL already exists.
	if _, ok := m.terse[*terse.ShortenedURL]; !ok {
		return ErrShortenedNotFound
	}

	// Update the Terse value in the Terse map.
	m.terse[*terse.ShortenedURL] = terse

	return nil
}

// UpsertTerse upserts the Terse into the TerseStore. This implementation has no network activity and ignores the given
// context.
func (m *MemTerse) UpsertTerse(_ context.Context, terse *models.Terse) (err error) {

	// Lock the Terse map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Upsert the Terse into the Terse map.
	m.terse[*terse.ShortenedURL] = terse

	return nil
}

// VisitsStore returns the VisitsStore.
func (m *MemTerse) VisitsStore() VisitsStore {
	return m.visitsStore
}
