package endpoints

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations/public"
	"github.com/MicahParks/terse-URL/storage"
)

// HandleRedirect creates and /{shortened} endpoint handler via a closure. It can perform redirects based on the
// shortened URL's Terse data. It will add visits to the VisitStore, if it exists.
func HandleRedirect(logger *zap.SugaredLogger, tmpl *template.Template, terseStore storage.TerseStore) public.TerseRedirectHandlerFunc {
	return func(params public.TerseRedirectParams) middleware.Responder {

		// Debug info.
		logger.Debugw("Parameters",
			"shortened", params.Shortened,
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the current time in the desired format.
		visitTime := strfmt.DateTime(time.Now())

		// Create the visit to represent this request.
		visit := &models.Visit{
			Accessed: &visitTime,
			Headers:  params.HTTPRequest.Header,
			IP:       &params.HTTPRequest.RemoteAddr, // TODO Use X-Forwarded-For if configured to do so.
		}

		// Get the Terse from the TerseStore.
		terse, err := terseStore.Read(ctx, params.Shortened, visit)
		if err != nil {

			// Log at the appropriate level.
			if errors.Is(err, storage.ErrShortenedNotFound) {
				logger.Infow("Shortened URL not found.",
					"shortened", params.Shortened,
					"error", err.Error(),
				)
			} else {
				logger.Errorw("Failed to get original URL from shortened.",
					"shortened", params.Shortened,
					"error", err.Error(),
				)
			}

			// Report the error to the client.
			return &public.TerseRedirectNotFound{}
		}

		// Check to see if an HTML file should be returned instead.
		if terse.JavascriptTracking || terse.MediaPreview != nil {

			// Create a buffer to write the populated HTML template with.
			buf := bytes.NewBuffer(nil)

			// If there is no error in populating the HTML template, return an HTML document to the client.
			if err = tmpl.Execute(buf, terse.MediaPreview); err == nil {
				return &public.TerseRedirectOK{Payload: ioutil.NopCloser(buf)}
			}

			// Failed to execute HTML template. Log the event. Reassign the error to nil. Attempt to issue standard
			// redirect.
			logger.Warnw("Failed to execute template. Attempting to perform standard redirect.",
				"shortened", params.Shortened,
				"error", err.Error(),
			) // TODO Different level?
			err = nil
		}

		// If the Terse data has an original URL, issue a standard redirect.
		if terse.OriginalURL != nil {
			return &public.TerseRedirectFound{
				Location: *terse.OriginalURL,
			}
		}

		// Log the event.
		logger.Warnw("Terse data did not contain original URL. Returning 404.",
			"shortened", params.Shortened,
		) // TODO Different level?
		return &public.TerseRedirectNotFound{}
	}
}
