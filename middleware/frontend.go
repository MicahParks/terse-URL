package middleware

import (
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/MicahParks/terseurl"
)

const (

	// frontendPrefix is the URL prefix used to access frontend assets.
	frontendPrefix = "/frontend/"
)

// FrontendMiddleware is the middleware used to server frontend assets.
func FrontendMiddleware(next http.Handler) (handler http.HandlerFunc, err error) {

	// Get the file system configuration from the environment.
	//
	// TODO Add function in configure package for grabbing this string value.
	frontendDir := os.Getenv("FRONTEND_STATIC_DIR")

	// Get the appropriate file system for the frontend assets.
	var fileSystem fs.FS
	if fileSystem, err = terseurl.FrontendFS(frontendDir); err != nil {
		return nil, err
	}

	// Create the file server.
	fileServer := http.FileServer(http.FS(fileSystem))

	// Create the HTTP handler via a closure.
	return func(writer http.ResponseWriter, request *http.Request) {

		// Redirect for the root file.
		if request.URL.Path == "/" || request.URL.Path == "/index.html" || request.URL.Path == strings.TrimSuffix(frontendPrefix, "/") {

			// Permanent redirect.
			http.Redirect(writer, request, frontendPrefix, http.StatusMovedPermanently)

			// Handle requests with the /frontend prefix.
		} else if strings.HasPrefix(request.URL.Path, frontendPrefix) {

			// Trim the prefix from the path so the file server behaves as expected.
			request.URL.Path = strings.TrimPrefix(request.URL.Path, frontendPrefix)

			// Serve the file system via HTTP.
			fileServer.ServeHTTP(writer, request)
		} else {

			// Follow the HTTP middleware pattern.
			next.ServeHTTP(writer, request)
		}
	}, nil
}
