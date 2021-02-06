package meta

import (
	"net/http"

	"github.com/MicahParks/terseurl/models"
)

// GetMeta performs an HTTP get on the given URL and returns the relevant HTML meta.
func GetMeta(u string) (og models.OpenGraph, twitter models.Twitter, err error) {

	// Perform an HTTP get on the URL.
	var resp *http.Response
	if resp, err = http.Get(u); err != nil { // TODO Use a non default HTTP client?
		return nil, nil, err
	}
	defer resp.Body.Close() // Ignore any error.

	// Return the relevant HTML meta.
	return previewTagInfo(resp.Body)
}
