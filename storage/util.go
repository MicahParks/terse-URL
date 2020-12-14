package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/MicahParks/ctxerrgroup"
	"go.etcd.io/bbolt"
)

const (

	// storageBbolt is the constant used when describing a storage backend as a bbolt file.
	storageBbolt = "bbolt"

	// storageMemory is the constant used when describing a storage backend only in memory.
	storageMemory = "memory"

	// storageNil is the constant used when describing a non-existent storage backend.
	storageNil = "nil"
)

var (

	// ErrShortenedNotFound indicates the given shortened URL was not found in the underlying storage.
	ErrShortenedNotFound = errors.New("the shortened URL was not found")

	// ErrShortenedExists indicates that an attempt was made to add a shortened URL that already existed.
	ErrShortenedExists = errors.New("the shortened URL already exists")

	// bboltTerseBucket is the bbolt bucket to use for Terse.
	bboltTerseBucket = []byte("terse")

	// bboltVisitsBucket is the bbolt bucket to use for Visits.
	bboltVisitsBucket = []byte("terseVisits")
)

type configuration struct {
	Type      string `json:"type"`
	BboltPath string `json:"bboltPath"`
}

type ctxCreator func() (ctx context.Context, cancel context.CancelFunc)

// TODO
func NewTerseStore(configJSON json.RawMessage, createCtx ctxCreator, errChan chan<- error, group *ctxerrgroup.Group, visitsStore VisitsStore) (terseStore TerseStore, storeType string, err error) {

	// Create the configuration.
	config := &configuration{}

	// If no JSON was given, use an in memory implementation.
	if len(configJSON) == 0 {
		config.Type = storageMemory
	} else {

		// Turn the configuration JSON into a Go structure.
		if err = json.Unmarshal(configJSON, config); err != nil {
			return nil, "", err
		}
	}

	// Create the appropriate TerseStore.
	switch config.Type {

	// Open a file as a bbolt database for the TerseStore.
	case storageBbolt:
		var db *bbolt.DB
		if db, err = openBbolt(config.BboltPath); err != nil {
			return nil, "", err
		}
		terseStore = NewBboltTerse(db, createCtx, group, bboltTerseBucket, visitsStore)

	// Use and in memory implementation of the VisitsStore by default.
	default:
		config.Type = storageMemory
		terseStore = NewMemTerse(createCtx, errChan, group, visitsStore)
	}

	return terseStore, config.Type, nil
}

// TODO
func NewVisitsStore(configJSON json.RawMessage) (visitsStore VisitsStore, storeType string, err error) {

	// Create the configuration.
	config := &configuration{}

	// If no configuration was give, return a nil VisitsStore.
	if len(configJSON) == 0 {
		return nil, storageNil, nil
	}

	// Turn the configuration JSON into a Go structure.
	if err = json.Unmarshal(configJSON, config); err != nil {
		return nil, "", err
	}

	// Create the appropriate VisitsStore.
	switch config.Type {

	// Open a file as a bbolt database for the VisitsStore.
	case storageBbolt:
		var db *bbolt.DB
		if db, err = openBbolt(config.BboltPath); err != nil {
			return nil, "", err
		}
		visitsStore = NewBboltVisits(db, bboltVisitsBucket)

	// Use and in memory implementation of the VisitsStore by default.
	default:
		config.Type = storageNil
		visitsStore = nil
	}

	return visitsStore, config.Type, nil
}

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

// TODO
func openBbolt(filePath string) (db *bbolt.DB, err error) {
	return bbolt.Open(filePath, 0666, nil)
}
