package endpoints

import (
	"fmt"
	"os"
	"path"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/restapi/operations"
)

func HandleFrontend(frontendDir string, logger *zap.SugaredLogger) operations.FrontendHandlerFunc {
	return func(params operations.FrontendParams) middleware.Responder {

		// Create the file path.
		filePath := path.Join(frontendDir, params.Path)

		// Stat the file real quick to make sure it exists.
		var err error
		if _, err = os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				logger.Infow("Failed to find requested file.",
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
		var file *os.File
		if file, err = os.Open(filePath); err != nil {
			logger.Errorw("Failed to open requested file.",
				"filePath", filePath,
				"error", err.Error(),
			)
		}

		// TODO Handle templating.

		return &operations.FrontendOK{
			Payload: file,
		}
	}
}
