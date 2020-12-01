package storage

import (
	"context"
	"errors"
	"time"
)

var (

	// ErrShortenedNotFound indicates the given shortened URL was not found in the underlying storage.
	ErrShortenedNotFound = errors.New("the shortened URL was not found")
)

type ctxCreator func() (ctx context.Context, cancel context.CancelFunc)

// deleteTerseBlocking deletes the Terse associated with the given shortened URL at the given deletion time. If an error
// occurs, it will be reported via errChan, if it exists.
func deleteTerseBlocking(createCtx ctxCreator, deleteAt time.Time, errChan chan<- error, shortened string, terseStore TerseStore) {

	// Figure out when the Terse needs to be deleted.
	duration := time.Until(deleteAt)

	// Wait until it is time to delete the Terse.
	if duration > 0 {
		time.Sleep(duration)
	}

	// Create a context for the deletion.
	ctx, cancel := createCtx()
	defer cancel()

	// Delete the Terse associated with the shortened URL.
	if err := terseStore.DeleteTerse(ctx, shortened); err != nil {

		// If the error channel exists, use it to report the error.
		if errChan != nil {
			errChan <- err
		}
	}
}
