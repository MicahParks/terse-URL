package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
	"github.com/MicahParks/terse-URL/storage"
)

func HandleDump(logger *zap.SugaredLogger, terseStore storage.TerseStore) operations.TerseDumpHandlerFunc {
	return func(params operations.TerseDumpParams) middleware.Responder {

		// Debug info.
		logger.Debug("Performing data dump.")

		// TODO Non-debug level log?

		// Create a request context.
		//
		// Maybe make a longer context if timing out.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the data dump.
		dump, err := terseStore.DumpAll(ctx)
		if err != nil {

			// Log at the appropriate level.
			logger.Errorw("Failed to perform data dump.",
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			message := "Failed to perform data dump."
			return &operations.TerseDumpDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		return &operations.TerseDumpOK{
			Payload: dump,
		}
	}
}
