package storage

import (
	"context"
	"errors"

	"github.com/MicahParks/terseurl/models"
)

var (

	// ErrKeyNotFound indicates that a given key was not found.
	ErrKeyNotFound = errors.New("the given key was not found") // TODO Remove ErrShortenedURLNotFound replace with this. Use fmt.Errorf to add detail.
)

// TODO Use a user's key as a UUID or string?

// AuthorizationStore is the Authorization data storage interface. It allows for Authorization data storage operations
// without needing to know of the Authorization data are stored.
//
// The underlying data structure, data structure 1, can be represented by a map[string]map[string]Authorization. Where
// the first string key is the user's unique identifier and the second string key is the shortened URL.
//
// There will be a second data structure, data structure 2, that will be map[string]map[string]struct{} which will be
// updated in tandem. The first string key will be the shortened URL and the second string key will be the user's unique
// identifier. The purpose of this second data type is for quick lookups relationships between shortened URLs and the
// users authorized to perform operations on them. This data structure should not be saved to persistent storage and
// should be built from the first data structure on creation. This is unless performance concerns arise from this
// strategy.
//
// The Authorization data will only reside in data structure 1.
//
// If a user is associated to a shortened URL at all it is assumed they have the ability to perform read operations on
// the associated shortened URL's Terse data.
//
// TODO Rename AuthorizationManager. Turn data structure 1 & 2 into an interface.
type AuthorizationStore interface {

	// Append upserts the given Authorization data. It works in an append first, overwrite second fashion. For example,
	// if the existing Authorization data indicates that a user currently has the shortened URLs X and Y, but the given
	// Authorization data indicates that the same user has the shortened URLs Y and Z, then the resulting Authorization
	// data for that user will include X, Y, and Z. Where Y's Authorization data is the most recently received.
	//
	// Both data structure 1 and data structure 2 should be updated during this operation.
	Append(ctx context.Context, usersShortened map[string]UserAuth) (err error)

	// Close closes the connection to the underlying storage.
	Close(ctx context.Context) (err error)

	// DeleteShortened deletes the Authorization data for the given shortened URLs for all users. If shortenedURLs is
	// nil or empty, all Authorization data are deleted. No error should be returned if a shortened URL is not found.
	//
	// This should first interact with data structure 2. During this interaction the affected users should be noted and
	// used to update the underlying data structure, data structure 1, afterwards.
	DeleteShortened(ctx context.Context, shortenedURLs []string) (err error)

	// DeleteUsers deletes the Authorization data for the given users. If users is nil or empty, all Authorization data
	// are deleted. No error should be returned if a user is not found.
	//
	// This should first interact directly with the underlying data structure, data structure 1. During this interaction
	// the affected shortened URLs should be noted and used to update data structure 2 afterwards.
	DeleteUsers(ctx context.Context, users []string) (err error)

	// Overwrite upserts the given Authorization data. It works in an overwrite only fashion. For example, if the
	// existing Authorization data indicates that a user currently has the shortened URLs X and Y, but the given
	// Authorization data indicates that the same user has the shortened URLs Y and Z, then the resulting Authorization
	// data for that user will only include Y and Z. Where Y's Authorization data is the most recently received.
	//
	// Both data structure 1 and data structure 2 should be updated during this operation.
	Overwrite(ctx context.Context, usersShortened map[string]UserAuth) (err error)

	// ReadUsers exports the Authorization data for the given users. If users is nil or empty, then all users'
	// Authorization data is expected.
	//
	// This should interact directly with the underlying data structure, data structure 1.
	ReadUsers(ctx context.Context, users []string) (usersShortened map[string]UserAuth, err error)

	// ReadShortened exports the Authorization data for the given shortened URLs. If shortenedURLs is nil or empty, then
	// all shortened URL Authorization data is expected.
	//
	// This should first interact with data structure 2 for faster lookups, then gather the Authorization data from data
	// structure 1.
	ReadShortened(ctx context.Context, shortenedURLs []string) (shortenedUserSet map[string]ShortenedAuth, err error)
}

// SummaryStore is the Summary data storage interface. It allows for Summary data storage operations without needing to
// know how the Summary data are stored. Summary data should not persist through a service restart.
type SummaryStore interface {

	// Close closes the connection to the underlying storage.
	Close(ctx context.Context) (err error)

	// Delete deletes the summary information for the given shortened URLs. If shortenedURLs is nil or empty, all
	// Summary data are deleted. No error should be returned if a shortened URL is not found.
	Delete(ctx context.Context, shortenedURLs []string) (err error)

	// IncrementVisitCount increments the visit count for the given shortened URL. It is called in separate goroutine.
	// The error must be storage.ErrShortenedNotFound if the shortened URL is not found.
	IncrementVisitCount(ctx context.Context, shortened string) (err error)

	// Read provides the summary information for the given shortened URLs. If shortenedURLs is nil or empty, all
	// summaries are returned. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
	Read(ctx context.Context, shortenedURLs []string) (summaries map[string]*models.Summary, err error)

	// Upsert upserts the summary information for the given shortened URL.
	Upsert(ctx context.Context, summaries map[string]*models.Summary) (err error)
}

// TerseStore is the Terse storage interface. It allows for Terse storage operations without needing to know how
// the Terse data are stored.
type TerseStore interface {

	// Close closes the connection to the underlying storage.
	Close(ctx context.Context) (err error)

	// Delete deletes the Terse data for the given shortened URLs. If shortenedURLs is nil or empty, all shortened URL
	// Terse data are deleted. There should be no error if a shortened URL is not found.
	Delete(ctx context.Context, shortenedURLs []string) (err error)

	// Read returns a map of shortened URLs to Terse data. If shortenedURLs is nil or empty, all shortened URL Terse
	// data are expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
	Read(ctx context.Context, shortenedURLs []string) (terseData map[string]*models.Terse, err error)

	// Summary summarizes the Terse data for the given shortened URLs. If shortenedURLs is nil or empty, then all
	// shortened URL Summary data are expected.
	Summary(ctx context.Context, shortenedURLs []string) (summaries map[string]*models.TerseSummary, err error)

	// Write writes the given Terse data according to the given operation. The error must be storage.ErrShortenedExists
	// if an Insert operation cannot be performed due to the Terse data already existing. The error must be
	// storage.ErrShortenedNotFound if an Update operation cannot be performed due to the Terse data not existing.
	Write(ctx context.Context, terseData map[string]*models.Terse, operation WriteOperation) (err error)
}

// VisitsStore is the Visits storage interface. It allows for Visits storage operations without needing to know how the
// Visits data are stored.
type VisitsStore interface {

	// Close closes the connection to the underlying storage.
	Close(ctx context.Context) (err error)

	// Delete deletes Visits data for the given shortened URLs. If shortenedURLs is nil or empty, then all Visits data
	// are deleted. No error should be given if a shortened URL is not found.
	Delete(ctx context.Context, shortenedURLs []string) (err error)

	// Insert inserts the given Visits data. The visits do not need to be unique, so the Visits data should be appended
	// to the data structure in storage.
	Insert(ctx context.Context, visitsData map[string][]models.Visit) (err error)

	// Read exports the Visits data for the given shortened URLs. If shortenedURLs is nil or empty, then all shortened
	// URL Visits data are expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
	Read(ctx context.Context, shortenedURLs []string) (visitsData map[string][]models.Visit, err error)

	// Summary summarizes the Visits data for the given shortened URLs. If shortenedURLs is nil or empty, then all
	// shortened URL Summary data are expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not
	// found.
	Summary(ctx context.Context, shortenedURLs []string) (summaries map[string]*models.VisitsSummary, err error)
}
