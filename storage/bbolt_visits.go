package storage

import (
	"context"
	"encoding/json"

	"go.etcd.io/bbolt"

	"github.com/MicahParks/terse-URL/models"
)

type BboltVisits struct {
	bbolt        *bbolt.DB
	visitsBucket []byte
}

func (b *BboltVisits) AddVisit(_ context.Context, shortened string, visit *models.Visit) (err error) {

	// Get the existing visits.
	var visits []*models.Visit
	if visits, err = b.readVisits(shortened); err != nil {
		return err
	}

	// Add the visits to the existing visits.
	visits = append(visits, visit)

	// Turn the visits into JSON data.
	var data []byte
	if data, err = json.Marshal(visits); err != nil {
		return err
	}

	// Open the bbolt database for writing, batch if possible.
	if err = b.bbolt.Batch(func(tx *bbolt.Tx) error {

		// Put the updated JSON data into the bucket.
		if err = tx.Bucket(b.visitsBucket).Put([]byte(shortened), data); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (b *BboltVisits) Close(_ context.Context) (err error) {
	return b.bbolt.Close()
}

func (b *BboltVisits) DeleteVisits(_ context.Context, shortened string) (err error) {

	// Open the bbolt database for writing, batch if possible.
	if err = b.bbolt.Batch(func(tx *bbolt.Tx) error {

		// Delete all of this shortened URL's visits from the bucket.
		if err = tx.Bucket(b.visitsBucket).Delete([]byte(shortened)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (b *BboltVisits) ReadVisits(_ context.Context, shortened string) (visits []*models.Visit, err error) {
	return b.readVisits(shortened)
}

func (b *BboltVisits) readVisits(shortened string) (visits []*models.Visit, err error) {

	// Open the bbolt database for reading.
	var data []byte
	if err = b.bbolt.View(func(tx *bbolt.Tx) error {

		// Get the Visits from the bucket.
		data = tx.Bucket(b.visitsBucket).Get([]byte(shortened))

		return nil
	}); err != nil {
		return nil, err
	}

	// Turn the JSON data into the Go structure.
	if err = json.Unmarshal(data, &visits); err != nil {
		return nil, err
	}

	return visits, nil
}
