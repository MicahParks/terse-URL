package auth

import (
	"errors"
)

var (

	// ErrKeyNotFound indicates that a given key was not found. This is indicative of a programming mistake.
	ErrKeyNotFound = errors.New("the given key was not found. Indicative of a logic error")
)

// Authorization represents Authorization data.
//
// If a user has any Authorization data associate with a shortened URL, it is implied that that user can read the
// shortened URL's Terse data.
//
// The owner has unrestricted access to a shortened URL's data and is the only one allowed to write Summary data and
// Visits data.
//
// The owner is the only one allowed to change other user's Authorization data for a given shortened URL.
//
// TODO Should this be apart of the swagger spec and communicated to the frontend to disable buttons and such?
type Authorization struct {
	Owner       bool `json:"owner"`
	ReadSummary bool `json:"read_summary"`
	ReadVisits  bool `json:"read_visits"`
	WriteTerse  bool `json:"write_terse"`
}
