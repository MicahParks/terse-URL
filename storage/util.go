package storage

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"

	"go.etcd.io/bbolt"

	"github.com/MicahParks/terseurl/models"
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

// configuration represents the configuration data gathered from the user to create a storage backend.
type configuration struct {
	Type      string `json:"type"`
	BboltPath string `json:"bboltPath"`
}

// CtxCreator is a function signature that creates a context and its cancel function.
type CtxCreator func() (ctx context.Context, cancel context.CancelFunc)

// NewSummaryStore creates a new SummaryStore from the given configJSON. The storeType return value is used for logging.
func NewSummaryStore(configJSON json.RawMessage) (summaryStore SummaryStore, storeType string, err error) {

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

	// Create the appropriate SummaryStore.
	switch config.Type {

	// Use and in memory implementation of the SummaryStore by default.
	default:
		config.Type = storageMemory
		summaryStore = NewMemSummary()
	}

	return summaryStore, config.Type, nil
}

// NewTerseStore creates a new TerseStore from the given configJSON. The storeType return value is used for logging.
func NewTerseStore(configJSON json.RawMessage) (terseStore TerseStore, storeType string, err error) {

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

		// Open the bbolt database file.
		var db *bbolt.DB
		if db, err = openBbolt(config.BboltPath); err != nil {
			return nil, "", err
		}

		// Create the bucket.
		if err = createBucket(db, bboltTerseBucket); err != nil {
			return nil, "", err
		}

		// Assign the interface implementation.
		terseStore = NewBboltTerse(db, bboltTerseBucket)

	// Use and in memory implementation of the TerseStore by default.
	default:
		config.Type = storageMemory
		terseStore = NewMemTerse()
	}

	return terseStore, config.Type, nil
}

// NewVisitsStore creates a new VisitsStore from the given configJSON. The storeType return value is used for logging.
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

	// Use and in memory implementation of the VisitsStore.
	case storageMemory:
		visitsStore = NewMemVisits()

	// Open a file as a bbolt database for the VisitsStore.
	case storageBbolt:

		// Open the bbolt database file.
		var db *bbolt.DB
		if db, err = openBbolt(config.BboltPath); err != nil {
			return nil, "", err
		}

		// Create the bucket.
		if err = createBucket(db, bboltVisitsBucket); err != nil {
			return nil, "", err
		}

		// Assign the interface implementation.
		visitsStore = NewBboltVisits(db, bboltVisitsBucket)

	// Use and in memory implementation of the VisitsStore by default.
	default:
		config.Type = storageNil
		visitsStore = nil
	}

	return visitsStore, config.Type, nil
}

// bytesToTerse transforms bytes to Terse data.
func bytesToTerse(data []byte) (terse models.Terse, err error) {
	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)
	if err = dec.Decode(&terse); err != nil {
		return models.Terse{}, err
	}
	return terse, nil
}

// bytesToVisits transforms bytes to Visits data.
func bytesToVisits(data []byte) (visits []models.Visit, err error) {
	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)
	if err = dec.Decode(&visits); err != nil {
		return nil, err
	}
	return visits, nil
}

// createBucket creates the given bucketName in the given bbolt database, if it doesn't already exist.
func createBucket(db *bbolt.DB, bucketName []byte) (err error) {
	if err = db.Update(func(tx *bbolt.Tx) error {
		if _, err = tx.CreateBucketIfNotExists(bucketName); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// openBbolt opens the file found at filePath as a bbolt database.
func openBbolt(filePath string) (db *bbolt.DB, err error) {
	return bbolt.Open(filePath, 0666, nil)
}

// summarizeTerse creates a *models.TerseSummary from a models.Terse.
func summarizeTerse(terse models.Terse) (summary *models.TerseSummary) {
	return &models.TerseSummary{
		OriginalURL:  terse.OriginalURL,
		RedirectType: terse.RedirectType,
		ShortenedURL: terse.ShortenedURL,
	}
}

// terseToBytes transforms Terse data to bytes.
func terseToBytes(terse models.Terse) (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err = enc.Encode(&terse); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// visitsToBytes transforms Visits data to bytes.
func visitsToBytes(visits []models.Visit) (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err = enc.Encode(visits); err != nil { // TODO Need to use *[]models.Visit?
		return nil, err
	}
	return buf.Bytes(), nil
}
