package endpoints

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
	"github.com/MicahParks/terse-URL/storage"

	"github.com/teris-io/shortid"
)

func HandleNew(invalidPaths []string, logger *zap.SugaredLogger, shortID *shortid.Shortid, shortIDParanoid bool, terseStore storage.TerseStore) operations.URLNewHandlerFunc {
	return func(params operations.URLNewParams, _ *models.JWTInfo) middleware.Responder {

		// Do not have debug level logging on in production, as it will log clog up the logs.
		logger.Debugw("Parameters",
			"terse", fmt.Sprintf("%+v", params.Terse),
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Confirm the original is a URL.
		var err error
		if _, err = url.Parse(*params.Terse.OriginalURL); err != nil {

			// Log with the appropriate level.
			message := `"originalURL" is not a properly formatted URL.`
			logger.Infow(message,
				"originalURL", *params.Terse.OriginalURL,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(400)
			return &operations.URLNewDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		// If no shortened URL was given, generate one.
		if params.Terse.ShortenedURL != "" {

			// Generate the shortened URL.
			var message string
			if params.Terse.ShortenedURL, message, err = generateRandom(ctx, logger, invalidPaths, params, shortID, shortIDParanoid, terseStore); err != nil { // TODO Verify reassignment to shortened.

				// Report the error to the client.
				code := int64(500)
				return &operations.URLNewDefault{Payload: &models.Error{
					Code:    &code,
					Message: &message,
				}}
			}
		}

		//// Change the deletion time to the desired format.
		//var deletionTime *time.Time
		//if !time.Time(params.Original.DeleteAt).IsZero() { // TODO Verify.
		//	formatted := time.Time(params.Original.DeleteAt)
		//	deletionTime = &formatted
		//}

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

		return &operations.URLNewOK{
			Payload: params.Terse.ShortenedURL,
		}
	}
}

// generateRandom will generate a random URL.
func generateRandom(ctx context.Context, logger *zap.SugaredLogger, invalidPaths []string, params operations.URLNewParams, shortID *shortid.Shortid, shortIDParanoid bool, terseStore storage.TerseStore) (shortened, message string, err error) {

	// Enter a loop to create a valid random shortened URL.
	for {

		// Generate a shortened URL that will redirect to the original.
		if shortened, err = shortID.Generate(); err != nil {

			// Log with the appropriate level.
			message = "Failed to generate random shortened URL."
			logger.Errorw(message,
				"deleteAt", params.Terse.DeleteAt,
				"originalURL", *params.Terse.OriginalURL,
				"error", err.Error(),
			)

			return shortened, message, err
		}

		// Confirm the randomly generated URL is not invalid.
		for _, invalid := range invalidPaths {
			if invalid == shortened {
				continue
			}
		}

		// Confirm the randomly generated URL is not already in use, if paranoid.
		if shortIDParanoid {
			if _, err = terseStore.GetTerse(ctx, shortened, nil, nil, context.Background()); err != nil {

				// If the shortened URL was not found, good.
				if errors.Is(err, storage.ErrShortenedNotFound) {
					err = nil
				} else {

					// Log with the appropriate level.
					message = "Failed to check if randomly generated URL is already in use."
					logger.Errorw(message,
						"error", err.Error(),
					)

					return shortened, message, err
				}
			}
		}

		// The URL is valid and not in use.
		break
	}

	return shortened, message, err
}
