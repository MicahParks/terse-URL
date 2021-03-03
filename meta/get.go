package meta

import (
	"bytes"
	"context"
	"net/http"

	"github.com/MicahParks/terseurl/models"
)

// Get performs an HTTP get on the given URL and returns the relevant HTML meta.
func Get(ctx context.Context, u string) (og models.OpenGraph, twitter models.Twitter, err error) {

	// Create the request.
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, u, bytes.NewReader(nil)); err != nil {
		return nil, nil, err
	}

	// TODO Add a user agent or something?

	// Perform an HTTP get on the URL.
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil { // TODO Use a non default HTTP client? Without following redirects?
		return nil, nil, err
	}
	defer resp.Body.Close() // Ignore any error.

	// Return the relevant HTML meta.
	return previewTagInfo(resp.Body)
}
