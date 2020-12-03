package endpoints

import (
	"context"
	"errors"
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

func HandleRandom(invalidPaths []string, logger *zap.SugaredLogger, shortID *shortid.Shortid, terseStore storage.TerseStore) operations.URLRandomHandlerFunc {
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
			message := `Parameter "original" is not a properly formatted URL.`
			logger.Infow(message,
				"original", params.Original.URL,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(400)
			return &operations.URLRandomDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		// Enter a loop to create a valid random shortened URL.
		var shortened string
		for {

			// Generate a shortened URL that will redirect to the original.
			if shortened, err = shortID.Generate(); err != nil {

				// Log with the appropriate level.
				message := "Failed to generate random shortened URL."
				logger.Errorw(message,
					"deleteAt", params.Original.DeleteAt,
					"original", params.Original.URL,
					"error", err.Error(),
				)

				// Report the error to the client.
				code := int64(500)
				return &operations.URLRandomDefault{Payload: &models.Error{
					Code:    &code,
					Message: &message,
				}}
			}

			// Confirm the randomly generated URL is not invalid.
			for _, invalid := range invalidPaths {
				if invalid == shortened {
					continue
				}
			}

			// Confirm the randomly generated URL is not already in use.
			if _, err = terseStore.GetTerse(ctx, shortened, nil, nil, context.Background()); err != nil {

				// If the shortened URL was not found, good.
				if errors.Is(err, storage.ErrShortenedNotFound) {
					err = nil
				} else {

					// Log with the appropriate level.
					message := "Failed to check if randomly generated URL is already in use."
					logger.Errorw(message,
						"error", err.Error(),
					)

					// Report the error to the client.
					code := int64(500)
					return &operations.URLRandomDefault{Payload: &models.Error{
						Code:    &code,
						Message: &message,
					}}
				}
			}

			// The URL is valid and not in use.
			break
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
			message := "Failed to upsert Terse pair into storage."
			logger.Errorw(message,
				"deleteAt", params.Original.DeleteAt,
				"original", params.Original.URL,
				"shortened", shortened,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
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
