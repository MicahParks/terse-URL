package auth

import (
	"context"

	"go.etcd.io/bbolt"
)

// BboltAuthorization TODO
type BboltAuthorization struct {
	db     *bbolt.DB
	bucket []byte
}

// NewBboltAuthorization TODO
func NewBboltAuthorization(db *bbolt.DB, authorizationBucket []byte) (authStore AuthorizationStore) {
	return BboltAuthorization{
		db:     db,
		bucket: authorizationBucket,
	}
}

func (b BboltAuthorization) Append(_ context.Context, usersShortened map[string]UserData) (err error) {
	panic("implement me")
}

func (b BboltAuthorization) Close(_ context.Context) (err error) {

	// Close the bbolt database file.
	return b.db.Close()
}

func (b BboltAuthorization) DeleteShortened(_ context.Context, shortenedURLs []string) (err error) {
	panic("implement me")
}

func (b BboltAuthorization) DeleteUsers(_ context.Context, users []string) (err error) {
	panic("implement me")
}

func (b BboltAuthorization) Overwrite(_ context.Context, usersShortened map[string]UserData) (err error) {
	panic("implement me")
}

func (b BboltAuthorization) ReadUsers(_ context.Context, users []string) (usersShortened map[string]UserData, err error) {
	panic("implement me")
}

func (b BboltAuthorization) ReadShortened(_ context.Context, shortenedURLs []string) (shortenedUsers map[string]ShortenedData, err error) {
	panic("implement me")
}
