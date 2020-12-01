package endpoints

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
	"github.com/MicahParks/terse-URL/storage"

	"github.com/teris-io/shortid"
)

func HandleRandom(logger *zap.SugaredLogger, shortID *shortid.Shortid, terseStore storage.TerseStore) operations.URLRandomHandlerFunc {
	return func(params operations.URLRandomParams, _ *models.JWTInfo) middleware.Responder {

		// Do not have debug level logging on in production, as it will log clog up the logs.
		logger.Debugw("Parameters",
			"original", fmt.Sprintf("%+v", params.Original),
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Confirm the original is a URL.
		var err error
		if _, err = url.Parse(params.Original.URL); err != nil {

			// Log with the appropriate level.
			logger.Infow("Client submitted original URL did not validate.",
				"original", params.Original.URL,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(400)
			message := `Parameter "original" is not a properly formatted URL.`
			return &operations.URLRandomDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		// Generate a shortened URL that will redirect to the original.
		var shortened string
		if shortened, err = shortID.Generate(); err != nil {

			// Log with the appropriate level.
			logger.Errorw("Failed to generate shortened URL.",
				"deleteAt", params.Original.DeleteAt,
				"original", params.Original.URL,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			message := "Failed to generate random shortened URL."
			return &operations.URLRandomDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		// Change the deletion time to the desired format.
		var deletionTime *time.Time
		if !time.Time(params.Original.DeleteAt).IsZero() { // TODO Verify.
			formatted := time.Time(params.Original.DeleteAt)
			deletionTime = &formatted
		}

		// Upsert the Terse pair into storage.
		if err = terseStore.UpsertTerse(ctx, deletionTime, params.Original.URL, shortened); err != nil {

			// Log with the appropriate level.
			logger.Errorw("Failed to upsert Terse pair.",
				"deleteAt", params.Original.DeleteAt,
				"original", params.Original.URL,
				"shortened", shortened,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			message := "Failed to upsert Terse pair into storage."
			return &operations.URLRandomDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		return &operations.URLRandomOK{
			Payload: shortened,
		}
	}
}
