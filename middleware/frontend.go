package middleware

import (
	"net/http"
	"strings"

	"github.com/MicahParks/terseurl"
)

const (

	// frontendPrefix is the URL prefix used to access frontend assets.
	frontendPrefix = "/frontend/"
)

// FrontendMiddleware is the middleware used to server frontend assets.
func FrontendMiddleware(frontendDir string, next http.Handler) (handler http.HandlerFunc, err error) {

	// Get the appropriate file system for the frontend assets.
	fileSystem, err := terseurl.FrontendFS(frontendDir)
	if err != nil {
		return nil, err
	}
	httpFileSystem := http.FS(fileSystem)

	// Create the HTTP handler via a closure.
	return func(writer http.ResponseWriter, request *http.Request) {

		// Redirect for the root file.
		if request.URL.Path == "/" || request.URL.Path == "/index.html" || request.URL.Path == strings.TrimPrefix(frontendPrefix, "/") {

			// Permanent redirect.
			http.Redirect(writer, request, frontendPrefix, 301)

			// Handle requests with the /frontend prefix.
		} else if strings.HasPrefix(request.URL.Path, frontendPrefix) {

			// Trim the prefix from the path so the file server behaves as expected.
			request.URL.Path = strings.TrimPrefix(request.URL.Path, frontendPrefix)

			// Serve the file system via HTTP.
			http.FileServer(httpFileSystem).ServeHTTP(writer, request)
		} else {

			// Follow the HTTP middleware pattern.
			next.ServeHTTP(writer, request)
		}
	}, nil
}
