package storage

import (
	"context"
	"sync"

	"github.com/MicahParks/terse-URL/models"
)

type MemVisits struct {
	mux    sync.RWMutex
	visits map[string][]*models.Visit
}

func NewMemVisits() (visitsStore VisitsStore) {
	return &MemVisits{
		visits: make(map[string][]*models.Visit),
	}
}

func (m *MemVisits) AddVisit(_ context.Context, shortened string, visit *models.Visit) (err error) {

	// Lock the Visits map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Confirm the shortened URL is a key in the map.
	var ok bool
	var visits []*models.Visit
	if visits, ok = m.visits[shortened]; !ok {
		visits = make([]*models.Visit, 0)
		m.visits[shortened] = visits
	}

	// Add the visits to the slice of visits for this shortened URL.
	m.visits[shortened] = append(visits, visit)

	return nil
}

func (m *MemVisits) Close(_ context.Context) (err error) {
	return nil
}

func (m *MemVisits) DeleteVisits(_ context.Context, shortened string) (err error) {

	// Lock the Visits map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Delete all visits for the shortened URL.
	delete(m.visits, shortened)

	return nil
}

func (m *MemVisits) ReadVisits(_ context.Context, shortened string) (visits []*models.Visit, err error) {

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
