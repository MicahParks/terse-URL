package storage

import (
	"context"
	"sync"

	"github.com/MicahParks/terse-URL/models"
)

type MemVisits struct {
	mux    *sync.Mutex
	visits map[string]*Visits
}

func NewMemVisits() (visitsStore VisitsStore) {
	return &MemVisits{
		mux:    &sync.Mutex{},
		visits: make(map[string]*Visits),
	}
}

func (m *MemVisits) Close(_ context.Context) (err error) {
	return nil
}

func (m *MemVisits) DeleteVisits(_ context.Context, shortened string) (err error) {

	// Lock the Visits map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Delete all visits.
	delete(m.visits, shortened)

	return nil
}

func (m *MemVisits) GetVisits(_ context.Context, shortened string) (visits []*models.Visit, err error) {

	// Lock the Visits map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Confirm the shortened URL is a key in the map.
	visitsStruct, ok := m.visits[shortened]
	if !ok {
		return nil, ErrShortenedNotFound
	}

	// Return the visits.
	return visitsStruct.Visits, nil
}

func (m *MemVisits) AddVisit(_ context.Context, shortened string, visit *models.Visit) (err error) {

	// Lock the Visits map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Confirm the shortened URL is a key in the map.
	visitsStruct, ok := m.visits[shortened]
	if !ok {
		return ErrShortenedNotFound
	}

	// Add the visits to the slice of visits for this shortened URL.
	visitsStruct.Visits = append(visitsStruct.Visits, visit)

	return nil
}
