package storage

import (
	"context"
	"fmt"

	"github.com/MicahParks/ctxerrgroup"
)

// TODO Remove VisitsStore and SummaryStore from TerseStore implementations.

type StoreManager struct { // TODO Rename.
	group        ctxerrgroup.Group
	summaryStore SummaryStore
	terseStore   TerseStore
	visitsStore  VisitsStore
}

func (s StoreManager) Close(ctx context.Context) (err error) {

	// Kill the worker pool.
	s.group.Kill()

	// TODO See if you can stick more than one error in with %w somehow...

	// Close the SummaryStore.
	var closeErr error
	if closeErr = s.SummaryStore().Close(ctx); closeErr != nil {
		err = fmt.Errorf("%v SummaryStore: %v", err, closeErr)
	}

	// Close the TerseStore.
	if closeErr = s.TerseStore().Close(ctx); closeErr != nil {
		err = fmt.Errorf("%v TerseStore: %v", err, closeErr)
	}

	// Close the VisitsStore.
	if closeErr = s.VisitsStore().Close(ctx); closeErr != nil {
		err = fmt.Errorf("%v VisitsStore: %v", err, closeErr)
	}

	return err
}

// SummaryStore returns the current SummaryStore.
func (s StoreManager) SummaryStore() (summaryStore SummaryStore) {
	return s.summaryStore
}

// TerseStore returns the current TerseStore.
func (s StoreManager) TerseStore() (terseStore TerseStore) {
	return s.terseStore
}

// VisitsStore returns the current VisitsStore.
func (s StoreManager) VisitsStore() (visitsStore VisitsStore) {
	return s.visitsStore
}
