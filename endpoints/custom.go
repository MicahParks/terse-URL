package endpoints

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
	"github.com/MicahParks/terse-URL/storage"
)

func HandleCustom(invalidPaths []string, logger *zap.SugaredLogger, terseStore storage.TerseStore) operations.URLCustomHandlerFunc {
	return func(params operations.URLCustomParams, _ *models.JWTInfo) middleware.Responder {

		// Do not have debug level logging on in production, as it will log clog up the logs.
		logger.Debugw("Parameters",
			"TersePair", fmt.Sprintf("%+v", params.TersePair),
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Confirm the original is a URL.
		var err error
		if _, err = url.Parse(params.TersePair.OriginalURL); err != nil {

			// Log with the appropriate level.
			message := `Parameter "originalURL" is not a properly formatted URL.`
			logger.Infow(message,
				"original", params.TersePair.OriginalURL,
				"error", err.Error())

			// Report the error to the client.
			code := int64(400)
			return &operations.URLCustomDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		// Confirm the custom shortened URL is URL safe.
		if _, err = url.Parse(params.TersePair.ShortenedURL); err != nil {

			// Log with the appropriate level.
			message := `Parameter "shortenedURL" is not a URL safe.`
			logger.Infow(message,
				"shortened", params.TersePair.ShortenedURL,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(400)
			return &operations.URLCustomDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		// Confirm the custom URL is not one of the given invalid URLs.
		for _, u := range invalidPaths {

			// If the custom URL is one of the given invalid URLs, report this to the client.
			if u == params.TersePair.ShortenedURL || strings.HasPrefix(params.TersePair.ShortenedURL, u+"/") { // TODO Need or condition?

				// Log with the appropriate level.
				message := "Invalid shortened URL"
				logger.Infow(message,
					"URL", params.TersePair.ShortenedURL,
				)

				// Report the error to the client.
				code := int64(400)
				return &operations.URLCustomDefault{Payload: &models.Error{
					Code:    &code,
					Message: &message,
				}}
			}
		}

		// Change the deletion time to the desired format.
		var deletionTime *time.Time
		if !time.Time(params.TersePair.DeleteAt).IsZero() { // TODO Verify.
			formatted := time.Time(params.TersePair.DeleteAt)
			deletionTime = &formatted
		}

		// Upsert the Terse pair into storage.
		if err = terseStore.UpsertTerse(ctx, deletionTime, params.TersePair.OriginalURL, params.TersePair.ShortenedURL); err != nil {

			// Log the server side failure.
			message := "Failed to upsert Terse pair into storage."
			logger.Errorw(message,
				"deleteAt", params.TersePair.DeleteAt,
				"original", params.TersePair.OriginalURL,
				"shortened", params.TersePair.ShortenedURL,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			return &operations.URLCustomDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		return &operations.URLCustomOK{
			Payload: params.TersePair.ShortenedURL, // TODO Query encode?
		}
	}
}
