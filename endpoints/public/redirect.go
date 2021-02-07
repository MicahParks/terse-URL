package public

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/meta"
	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/restapi/operations/public"
	"github.com/MicahParks/terseurl/storage"
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

		// TODO Validate OriginalURL, if needed. Like if empty.

		// Check to see if a 301 redirect needs to be issued.
		if terse.RedirectType == models.RedirectTypeNr301 {
			return &public.TerseRedirectMovedPermanently{Location: terse.OriginalURL}
		}

		// Check to see if an HTML file should be returned instead.
		if terse.MediaPreview != nil && terse.JavascriptTracking || terse.RedirectType == models.RedirectTypeJs || terse.RedirectType == models.RedirectTypeMeta { // TODO Verify logic behind this if statement.

			// Create a buffer to write the populated HTML template with.
			buf := bytes.NewBuffer(nil)

			// Create the proper metadata for the HTML page.
			previewMeta := meta.Preview{
				MediaPreview: *terse.MediaPreview,
				Redirect:     terse.OriginalURL,
				RedirectType: terse.RedirectType,
			}

			logger.Debugw("",
				"og", fmt.Sprintf("%+v", terse.MediaPreview.Og),
				"twitter", fmt.Sprintf("%+v", terse.MediaPreview.Twitter),
			)

			// If there is no error in populating the HTML template, return an HTML document to the client.
			if err = tmpl.Execute(buf, previewMeta); err == nil {
				return &public.TerseRedirectOK{Payload: ioutil.NopCloser(buf)}
			}

			// Failed to execute HTML template. Log the event. Reassign the error to nil. Perform the default redirect.
			logger.Warnw("Failed to execute template.",
				"shortened", params.Shortened,
				"error", err.Error(),
			)
			err = nil
		}

		// Check if a 302 was desired. If not, log.
		if terse.RedirectType != models.RedirectTypeNr302 {
			logger.Warnw("The desired type couldn't be formed. A 302 will be issued instead.",
				"redirectType", terse.RedirectType,
			)
		}

		// Issue a standard temporary redirect.
		return &public.TerseRedirectFound{
			Location: terse.OriginalURL,
		}
	}
}
