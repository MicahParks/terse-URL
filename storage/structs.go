package storage

import (
	"time"

	"github.com/MicahParks/terse-URL/models"
)

// Terse TODO
type Terse struct {
	DeleteAt  *time.Time `bson:"deleteAt"` // This is a pointer so it can be nil.
	Original  string     `bson:"original"`
	Shortened string     `bson:"_id"`
}

type Visits struct {
	Original  string          `bson:"original"`
	Shortened string          `bson:"_id"`
	Visits    []*models.Visit `bson:"visits"`
}
