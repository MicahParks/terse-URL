package storage

import (
	"context"
	"sync"

	"github.com/MicahParks/terse-URL/models"
)

// MemVisits is a VisitsStore implementation that stores all data in a Go map in memory.
type MemVisits struct {
	mux    sync.RWMutex
	visits map[string][]*models.Visit
}

// NewMemVisits creates a new MemVisits.
func NewMemVisits() (visitsStore VisitsStore) {
	return &MemVisits{
		visits: make(map[string][]*models.Visit),
	}
}

// Add adds the visit to the visits store. This implementation has no network activity and ignores the given
// context.
func (m *MemVisits) Add(_ context.Context, shortened string, visit *models.Visit) (err error) {

	// Lock the Visits map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Confirm the shortened URL is a key in the map.
	visits, ok := m.visits[shortened]
	if !ok {
		visits = make([]*models.Visit, 0)
		m.visits[shortened] = visits
	}

	// Add the visits to the slice of visits for this shortened URL.
	m.visits[shortened] = append(visits, visit)

	return nil
}

// Close lets the garbage collector take care of the old Visits data. This implementation has no network activity and
// ignores the given context.
func (m *MemVisits) Close(_ context.Context) (err error) {
	m.visits = make(map[string][]*models.Visit)
	return nil
}

// Delete deletes data according to the del argument. This implementation has no network activity and ignores the
// given context.
func (m *MemVisits) Delete(_ context.Context, del models.Delete) (err error) {

	// Confirm the deletion of Visits data.
	if del.Visits == nil || *del.Visits {

		// Lock the Visits map for async safe use.
		m.mux.Lock()

		// Reassign the Visits map and the garbage collector will take care of the old data.
		m.visits = make(map[string][]*models.Visit)

		// Unlock the Visits map. The write operation is over.
		m.mux.Unlock()
	}

	return nil
}

// DeleteOne deletes data according to the del argument for the shortened URL. No error will be given if the shortened
// URL is not found. This implementation has no network activity and ignores the given context.
func (m *MemVisits) DeleteOne(_ context.Context, del models.Delete, shortened string) (err error) {

	// Confirm the deletion of Visits data.
	if del.Visits == nil || *del.Visits {

		// Lock the Visits map for async safe use.
		m.mux.Lock()

		// Delete all visits for the shortened URL.
		delete(m.visits, shortened)

		// Unlock the Visits map. The write operation is over.
		m.mux.Unlock()
	}

	return nil
}

// Export exports all exports all visits data. This implementation has no network activity and ignores the given
// context.
func (m *MemVisits) Export(_ context.Context) (allVisits map[string][]*models.Visit, err error) {

	// Lock the Visits map for async safe use.
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.visits, nil
}

// ExportOne gets all visits to the shortened URL. The error storage.ErrShortenedNotFound will be given if the shortened
// URL is not found. This implementation has no network activity and ignores the given context.
func (m *MemVisits) ExportOne(_ context.Context, shortened string) (visits []*models.Visit, err error) {

	// Lock the Visits map for async safe use.
	m.mux.RLock()
	defer m.mux.RUnlock()

	// Confirm the shortened URL is a key in the map.
	var ok bool
	if visits, ok = m.visits[shortened]; !ok {
		return nil, ErrShortenedNotFound
	}

	// Return the visits.
	return visits, nil
}

// Import imports the given export's data. If del is not nil, data will be deleted accordingly. If del is nil, data
// may be overwritten, but unaffected data will be untouched. This implementation has no network activity and ignores
// the given context.
func (m *MemVisits) Import(ctx context.Context, del *models.Delete, export map[string]models.Export) (err error) {

	// Check if data needs to be deleted before importing.
	if del != nil {
		if err = m.Delete(ctx, *del); err != nil {
			return err
		}
	}

	// Lock the Visits map for async safe usage.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Write every shortened URL's Visits data to the Visits map.
	for shortened, exp := range export {
		m.visits[shortened] = exp.Visits
	}

	return nil
}
