package storage

import (
	"context"
	"sync"

	"github.com/MicahParks/ctxerrgroup"

	"github.com/MicahParks/terseurl/models"
)

// MemTerse is a TerseStore implementation that stores all data in a Go map in memory.
type MemTerse struct {
	createCtx    ctxCreator
	errChan      chan<- error
	group        *ctxerrgroup.Group
	mux          sync.RWMutex
	summaryStore SummaryStore
	terse        map[string]*models.Terse
	visitsStore  VisitsStore
}

// NewMemTerse creates a new MemTerse given the required assets.
func NewMemTerse(createCtx ctxCreator, errChan chan<- error, group *ctxerrgroup.Group, summaryStore SummaryStore, visitsStore VisitsStore) (terseStore TerseStore) {
	return &MemTerse{
		createCtx:    createCtx,
		errChan:      errChan,
		group:        group,
		terse:        make(map[string]*models.Terse),
		summaryStore: summaryStore,
		visitsStore:  visitsStore,
	}
}

// Close closes the connection to the underlying storage. The ctxerrgroup will be killed. This will not close the
// connection to the VisitsStore.
func (m *MemTerse) Close(_ context.Context) (err error) {

	// Kill the worker pool.
	m.group.Kill()

	// Let the garbage collector take care of the old Terse data.
	m.terse = make(map[string]*models.Terse)

	return nil
}

// CreateSummaryStore creates the SummaryStore based on the existing VisitsStore data.
func (m *MemTerse) CreateSummaryStore(ctx context.Context) (summaries map[string]models.TerseSummary, err error) {

	// Create the return map.
	summaries = make(map[string]models.TerseSummary)

	// Confirm the visitsStore is not nil.
	if m.visitsStore != nil {

		// Get the map of counts.
		var counts map[string]uint
		counts, err = m.visitsStore.ExportCounts(ctx)
		if err != nil {
			return nil, err
		}

		// Lock the Terse map for async safe use.
		m.mux.RLock()

		// Create the Terse summary data by combining Visits counts and Terse data.
		for shortened, count := range counts {
			terse := m.terse[shortened]
			summaries[shortened] = models.TerseSummary{
				OriginalURL:  terse.OriginalURL,
				RedirectType: terse.RedirectType,
				ShortenedURL: terse.ShortenedURL,
				VisitCount:   int64(count), // TODO Conversion may cause data loss.
			}
		}

		// Unlock the Terse map for async safe use.
		m.mux.RUnlock()
	}

	return summaries, nil
}

// Delete deletes data according to the del argument. If the VisitsStore is not nil, then the same method will be
// called for the associated VisitsStore.
func (m *MemTerse) Delete(ctx context.Context, del models.Delete) (err error) {

	// Delete Visits data if required.
	if del.Visits == nil || *del.Visits && m.visitsStore != nil {
		if err = m.visitsStore.Delete(ctx, del); err != nil {
			return err
		}
	}

	// Confirm the deletion of Terse data.
	if del.Terse == nil || *del.Terse {

		// Lock the Terse Map for async safe use.
		m.mux.Lock()

		// Reassign the Terse map and the garbage collector will take care of the old data.
		m.terse = make(map[string]*models.Terse)

		// Unlock the Terse map. The write operation is over.
		m.mux.Unlock()
	}

	return nil
}

// DeleteSome deletes data according to the del argument for the given shortened URL. No error should be given if
// the shortened URL is not found. If the VisitsStore is not nil, then the same method will be called for the
// associated VisitsStore.
func (m *MemTerse) DeleteSome(ctx context.Context, del models.Delete, shortenedURLs []string) (err error) {

	// Delete Visits data if required.
	if del.Visits == nil || *del.Visits && m.visitsStore != nil {
		if err = m.visitsStore.DeleteSome(ctx, del, shortenedURLs); err != nil {
			return err
		}
	}

	// Check to make sure if Terse data needs to be deleted.
	if del.Terse == nil || *del.Terse {

		// Lock the Terse map for async safe use.
		m.mux.Lock()

		// Iterate through the shortened URLs.
		for _, shortened := range shortenedURLs {

			// Delete the Terse from the Terse map.
			delete(m.terse, shortened)
		}

		// Unlock the Terse map. The write operation is over.
		m.mux.Unlock()
	}

	return nil
}

// Export returns a map of shortened URLs to export data.
func (m *MemTerse) Export(ctx context.Context) (export map[string]models.Export, err error) {

	// Create the export map.
	export = make(map[string]models.Export)

	// Create the slice of shortened URLs. Preallocate the memory for faster insertion.
	shortenedURLs := make([]string, len(m.terse))
	index := 0
	for shortened := range m.terse {
		shortenedURLs[index] = shortened
		index++
	}

	// Create the slice map of Visits.
	var visits map[string][]*models.Visit
	if visits, err = m.visitsStore.Export(ctx); err != nil {
		return nil, err
	}

	// Lock the Terse map for async safe use.
	m.mux.RLock()

	// Iterate through the shortened URLs and add them to the map.
	for shortened, terse := range m.terse {

		// Add the shortened URL and its visits to the data export.
		export[shortened] = models.Export{
			Terse:  terse,
			Visits: visits[shortened], // TODO Check if ok?
		}
	}

	// Unlock the Terse map for async safe use.
	m.mux.RUnlock()

	return export, nil
}

// ExportSome returns a export of Terse and Visit data for a given shortened URL. The error must be
// storage.ErrShortenedNotFound if the shortened URL is not found.
func (m *MemTerse) ExportSome(ctx context.Context, shortenedURLs []string) (export map[string]models.Export, err error) {

	// Create the return map.
	export = make(map[string]models.Export)

	// Create the slice map of Visits.
	var visits map[string][]*models.Visit
	if visits, err = m.visitsStore.ExportSome(ctx, shortenedURLs); err != nil {
		return nil, err
	}

	// Lock the Terse map for async safe use.
	m.mux.RLock()
	defer m.mux.RUnlock()

	// Iterate through the shortened URLs.
	for _, shortened := range shortenedURLs {

		// Get the Terse.
		terse, ok := m.terse[shortened]
		if !ok {
			return nil, ErrShortenedNotFound
		}

		// Combine the data and add to the return slice.
		export[shortened] = models.Export{
			Terse:  terse,
			Visits: visits[shortened], // TODO Check if ok?
		}
	}

	return export, nil
}

// ExportTerse exports all Terse data as a map of shortened URLs to Terse data.
func (m *MemTerse) ExportTerse(_ context.Context) (terse map[string]*models.Terse, err error) {

	// Lock the Terse map for async safe use.
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.terse, nil
}

// Import imports the given export's data. If del is not nil, data will be deleted accordingly. If del is nil, data
// may be overwritten, but unaffected data will be untouched. If the VisitsStore is not nil, then the same method
// will be called for the associated VisitsStore.
func (m *MemTerse) Import(ctx context.Context, del *models.Delete, export map[string]models.Export) (err error) {

	// Check if data needs to be deleted before importing.
	if del != nil {
		if err = m.Delete(ctx, *del); err != nil {
			return err
		}
	}

	// Import the Visits data.
	if m.visitsStore != nil {

		// Import the Visits data. Never pass in a deletion data structure, because that already happened above.
		if err = m.visitsStore.Import(ctx, nil, export); err != nil {
			return err
		}
	}

	// Lock the Terse map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Write every shortened URL's Terse data to the Terse map.
	for shortened, exp := range export {
		m.terse[shortened] = exp.Terse
	}

	return nil
}

// Insert adds a Terse to the TerseStore. The shortened URL will be active after this. The error will be
// storage.ErrShortenedExists if the shortened URL is already present.
func (m *MemTerse) Insert(_ context.Context, terse *models.Terse) (err error) {

	// Lock the Terse map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Check to see if the shortened URL already exists.
	if _, ok := m.terse[terse.ShortenedURL]; ok {
		return ErrShortenedExists
	}

	// Add the shortened URL to the Terse map.
	m.terse[terse.ShortenedURL] = terse

	return nil
}

// Read retrieves all non-Visit Terse data give its shortened URL. A nil visit may be passed in and the visit should
// not be recorded. The error must be storage.ErrShortenedNotFound if the shortened URL is not found.
func (m *MemTerse) Read(_ context.Context, shortened string, visit *models.Visit) (terse *models.Terse, err error) {

	// TODO Combine SummaryStore and VisitsStore work items?
	// TODO Figure out how to make implementation of these methods more generic. Maybe add some functions that accept TerseStore.

	// Only update the SummaryStore and VisitsStore if there wasn't an error.
	defer func() {
		if err == nil {

			// Increment the number of times the shortened URL has been visited. Do this in a separate goroutine so the response
			// is faster.
			if m.summaryStore != nil {
				ctx, cancel := m.createCtx()
				go m.group.AddWorkItem(ctx, cancel, func(workCtx context.Context) (err error) {
					return m.summaryStore.IncrementVisitCount(ctx, shortened)
				})
			}

			// Track the visit to this shortened URL. Do this in a separate goroutine so the response is faster.
			if visit != nil && m.visitsStore != nil {
				ctx, cancel := m.createCtx()
				go m.group.AddWorkItem(ctx, cancel, func(workCtx context.Context) (err error) {
					return m.visitsStore.Add(workCtx, shortened, visit)
				})
			}
		}
	}()

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

// SummaryStore returns the VisitsStore.
func (m *MemTerse) SummaryStore() SummaryStore {
	return m.summaryStore
}

// Update assumes the Terse already exists. It will override all of its values. The error must be
// storage.ErrShortenedNotFound if the shortened URL is not found.
func (m *MemTerse) Update(_ context.Context, terse *models.Terse) (err error) {

	// Lock the Terse map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Check to see if the shortened URL already exists.
	if _, ok := m.terse[terse.ShortenedURL]; !ok {
		return ErrShortenedNotFound
	}

	// Update the Terse value in the Terse map.
	m.terse[terse.ShortenedURL] = terse

	return nil
}

// Upsert will upsert the Terse into the backend storage.
func (m *MemTerse) Upsert(_ context.Context, terse *models.Terse) (err error) {

	// Lock the Terse map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Upsert the Terse into the Terse map.
	m.terse[terse.ShortenedURL] = terse

	return nil
}

// VisitsStore returns the VisitsStore.
func (m *MemTerse) VisitsStore() VisitsStore {
	return m.visitsStore
}
