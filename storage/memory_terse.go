package storage

import (
	"context"
	"sync"

	"github.com/MicahParks/ctxerrgroup"

	"github.com/MicahParks/terse-URL/models"
)

// TODO Function comments.

// TODO
type MemTerse struct {
	createCtx   ctxCreator
	errChan     chan<- error
	group       *ctxerrgroup.Group
	mux         sync.RWMutex
	terse       map[string]*models.Terse
	visitsStore VisitsStore
}

// TODO
func NewMemTerse(createCtx ctxCreator, errChan chan<- error, group *ctxerrgroup.Group, visitsStore VisitsStore) (terseStore TerseStore) {
	return &MemTerse{
		createCtx:   createCtx,
		errChan:     errChan,
		group:       group,
		terse:       make(map[string]*models.Terse),
		visitsStore: visitsStore,
	}
}

func (m *MemTerse) Close(_ context.Context) (err error) {
	return nil
}

func (m *MemTerse) CreateTerse(_ context.Context, terse *models.Terse) (err error) {

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

func (m *MemTerse) DeleteTerse(_ context.Context, shortened string) (err error) {

	// Lock the Terse map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Delete the Terse from the Terse map.
	delete(m.terse, shortened)

	return nil
}

func (m *MemTerse) Dump(ctx context.Context, shortened string) (dump models.Dump, err error) {

	// Lock the Terse map for async safe use.
	m.mux.RLock()
	defer m.mux.RUnlock()

	// Get the Terse.
	terse, ok := m.terse[shortened]
	if !ok {
		return models.Dump{}, ErrShortenedNotFound
	}

	// Get the visits for the Terse.
	var visits []*models.Visit
	if visits, err = m.visitsStore.ReadVisits(ctx, shortened); err != nil {
		return models.Dump{}, err
	}

	return models.Dump{
		Terse:  terse,
		Visits: visits,
	}, nil
}

func (m *MemTerse) DumpAll(ctx context.Context) (dump map[string]models.Dump, err error) {

	// Lock the Terse map for async safe use.
	m.mux.RLock()
	defer m.mux.RUnlock()

	// Create the dump map.
	dump = make(map[string]models.Dump)

	// Iterate through the shortened URLs and add them to the map.
	for shortened, terse := range m.terse {

		// Get the visits for the Terse.
		var visits []*models.Visit
		if visits, err = m.visitsStore.ReadVisits(ctx, *terse.ShortenedURL); err != nil {
			return nil, err
		}

		// Add the shortened URL and its visits to the data dump.
		dump[shortened] = models.Dump{
			Terse:  terse,
			Visits: visits,
		}
	}

	return dump, nil
}

func (m *MemTerse) GetTerse(_ context.Context, shortened string, visit *models.Visit) (terse *models.Terse, err error) {

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
	if ok {
		return nil, ErrShortenedExists
	}

	return terse, nil
}

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

func (m *MemTerse) UpsertTerse(_ context.Context, terse *models.Terse) (err error) {

	// Lock the Terse map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Upsert the Terse into the Terse map.
	m.terse[*terse.ShortenedURL] = terse

	return nil
}

func (m *MemTerse) VisitsStore() VisitsStore {
	return m.visitsStore
}
