package auth

import (
	"context"
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
	Append(ctx context.Context, usersShortened map[string]UserData) (err error)

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
	Overwrite(ctx context.Context, usersShortened map[string]UserData) (err error)

	// ReadUsers exports the Authorization data for the given users. If users is nil or empty, then all users'
	// Authorization data is expected.
	//
	// This should interact directly with the underlying data structure, data structure 1.
	ReadUsers(ctx context.Context, users []string) (usersShortened map[string]UserData, err error)

	// ReadShortened exports the Authorization data for the given shortened URLs. If shortenedURLs is nil or empty, then
	// all shortened URL Authorization data is expected.
	//
	// This should first interact with data structure 2 for faster lookups, then gather the Authorization data from data
	// structure 1.
	ReadShortened(ctx context.Context, shortenedURLs []string) (shortenedUserSet map[string]ShortenedData, err error)
}
