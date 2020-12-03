package endpoints

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
)

func HandleFrontend(frontendDir string, logger *zap.SugaredLogger) operations.FrontendHandlerFunc {
	return func(params operations.FrontendParams) middleware.Responder {

		// Do not have debug level logging on in production, as it will log clog up the logs.
		logger.Debugw("Parameters",
			"path", params.Path,
			"original", params.Original,
			"shortened", params.Shortened,
		)

		// Create the file path.
		filePath := path.Join(frontendDir, params.Path)

		// Stat the file real quick to make sure it exists.
		var err error
		if _, err = os.Stat(filePath); err != nil {

			// Log with the appropriate level.
			if os.IsNotExist(err) {
				logger.Warnw("Failed to find requested file.",
					"filePath", filePath,
				)
			} else {
				logger.Errorw("Failed to find requested file.",
					"filePath", filePath,
					"error", err.Error(),
				)
			}

			return &operations.FrontendNotFound{
				Payload: fmt.Sprintf(`File not found: "%s".`, params.Path),
			}
		}

		// Open the file for reading.
		var fileData []byte
		if fileData, err = ioutil.ReadFile(filePath); err != nil {
			logger.Errorw("Failed to open requested file.",
				"filePath", filePath,
				"error", err.Error(),
			)
		}

		// Template the values of the frontend HTML asset, if required.
		if err = buildHTML(fileData, params); err != nil {

			// Log with the appropriate level.
			message := "Failed to build HTML template."
			logger.Errorw(message,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			return &operations.FrontendDefault{
				Payload: &models.Error{
					Code:    &code,
					Message: &message,
				},
			}
		}

		return &operations.FrontendOK{
			Payload: ioutil.NopCloser(bytes.NewReader(fileData)), // File is already closed.
		}
	}
}
