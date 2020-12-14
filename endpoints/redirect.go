package endpoints

import (
	"errors"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations/public"
	"github.com/MicahParks/terse-URL/storage"
)

// HandleExport creates and /{shortened} endpoint handler via a closure. It can perform redirects based on the shortened
// URL's Terse data. It will add visits to the VisitStore, if it exists.
func HandleRedirect(logger *zap.SugaredLogger, terseStore storage.TerseStore) public.TerseRedirectHandlerFunc {
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
		terse, err := terseStore.ReadTerse(ctx, params.Shortened, visit)
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

		return &public.TerseRedirectFound{
			Location: *terse.OriginalURL,
		}
	}
}
