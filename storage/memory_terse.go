package storage

import (
	"context"
	"sync"
	"time"

	"gitlab.com/MicahParks/ctxerrgroup"

	"github.com/MicahParks/terse-URL/models"
)

type MemTerse struct {
	createCtx   ctxCreator
	errChan     chan<- error
	group       *ctxerrgroup.Group
	mux         *sync.Mutex
	terse       map[string]*Terse
	visitsStore VisitsStore
}

func NewMemTerse(createCtx ctxCreator, errChan chan<- error, group *ctxerrgroup.Group, visitsStore VisitsStore) (terseStore TerseStore) {
	return &MemTerse{
		createCtx:   createCtx,
		errChan:     errChan,
		group:       group,
		mux:         &sync.Mutex{},
		terse:       make(map[string]*Terse),
		visitsStore: visitsStore,
	}
}

func (m *MemTerse) Close(_ context.Context) (err error) {
	return nil
}

func (m *MemTerse) ScheduleDeletions(_ context.Context) (err error) {

	//// Lock the Terse map for async safe use.
	//m.mux.Lock()
	//defer m.mux.Unlock()
	//
	//// Iterate through all Terse.
	//for _, terse := range m.terse {
	//
	//	// If the Terse has a deletion time, schedule it asynchronously.
	//	if terse.DeleteAt != nil {
	//		go deleteTerseBlocking(*terse.DeleteAt, m.errChan, terse.Shortened, m)
	//	}
	//}

	// There is no use case for this function, since memory should not be reloaded on program restart.
	return nil
}

func (m *MemTerse) DeleteTerse(_ context.Context, shortened string) (err error) {

	// Lock the Terse map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Delete the Terse.
	delete(m.terse, shortened)

	return nil
}

func (m *MemTerse) GetTerse(_ context.Context, shortened string, visit *models.Visit, visitCancel context.CancelFunc, visitCtx context.Context) (original string, err error) {

	// Lock the Terse map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Track the visit to this shortened URL. Do this in a separate goroutine so the response is faster.
	if m.visitsStore != nil {
		go m.group.AddWorkItem(visitCtx, visitCancel, func(workCtx context.Context) (err error) {
			return m.visitsStore.AddVisit(workCtx, shortened, visit)
		})
	}

	// Confirm the shortened URL is a key in the map.
	terse, ok := m.terse[shortened]
	if !ok {
		return "", ErrShortenedNotFound
	}

	// Return the shortened URL's original URL.
	return terse.Original, nil
}

func (m *MemTerse) UpsertTerse(_ context.Context, deleteAt *time.Time, original, shortened string) (err error) {

	// Lock the Terse map for async safe use.
	m.mux.Lock()
	defer m.mux.Unlock()

	// Assign the shortened URL as the key and create a Terse for the map.
	m.terse[shortened] = &Terse{
		DeleteAt:  deleteAt,
		Original:  original,
		Shortened: shortened,
	}

	// Schedule the Terse for deletion, if a deletion time was given.
	if deleteAt != nil {
		go deleteTerseBlocking(m.createCtx, *deleteAt, m.errChan, shortened, m)
	}

	return nil
}

func (m *MemTerse) VisitsStore() VisitsStore {
	return m.visitsStore
}
