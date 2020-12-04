package endpoints

import (
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
	"github.com/MicahParks/terse-URL/storage"
)

func HandleGet(logger *zap.SugaredLogger, terseStore storage.TerseStore) operations.URLGetHandlerFunc {
	return func(params operations.URLGetParams) middleware.Responder {

		// Do not have debug level logging on in production, as it will log clog up the logs.
		logger.Debugw("Parameters",
			"shortened", fmt.Sprintf("%+v", params.Shortened),
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the current time of this request in the desire format.
		visitTime := strfmt.DateTime(time.Now())

		// Create the visit to represent this request.
		visit := &models.Visit{
			Accessed: &visitTime,
			IP:       &params.HTTPRequest.RemoteAddr,
			Headers:  params.HTTPRequest.Header,
		}

		// Create another context for the VisitStore interactions.
		visitCtx, visitCancel := configure.DefaultCtx()

		// Get the original URL from storage.
		var err error
		var original string
		if original, err = terseStore.GetTerse(ctx, params.Shortened, visit, visitCancel, visitCtx); err != nil {

			// Log with the appropriate level.
			logger.Warnw("Failed to find requested shortened URL.",
				"shortened", params.Shortened,
				"error", err.Error(),
			)

			// Assume the Terse pair was missing and return a 404 so web browsers behave normally.
			return &operations.URLGetNotFound{}
		}

		return &operations.URLGetFound{
			Location: original, // TODO Verify.
		}
	}
}
