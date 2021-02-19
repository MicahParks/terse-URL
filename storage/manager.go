package storage

import (
	"context"

	"github.com/MicahParks/ctxerrgroup"
)

// TODO Remove VisitsStore and SummaryStore from TerseStore implementations.

type StoreManager struct { // TODO Rename.
	group   ctxerrgroup.Group
	wrapped TerseStore
}

func (t TerseStoreWrapper) Close(ctx context.Context) (err error) {

	// Kill the worker pool.
	t.group.Kill()

	return t.wrapped.Close(ctx)
}
