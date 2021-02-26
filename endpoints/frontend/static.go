package frontend

import (
	"io/fs"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/restapi/operations/frontend"
)

// HandleFrontendStatic creates and /frontend/{fileName} endpoint handler via a closure. It servers static frontend
// assets.
func HandleFrontendStatic(logger *zap.SugaredLogger, fileSystem fs.FS) frontend.FrontendStaticHandlerFunc {
	return func(params frontend.FrontendStaticParams) middleware.Responder {

		// Debug info.
		logger.Debugw("Parameters",
			"fileName", params.FileName,
		)

		// Attempt to read the file from the given file system.
		file, err := fileSystem.Open(params.FileName)
		if err != nil {

			// Log at the appropriate level.
			logger.Warnw("Failed to open requested file.",
				"fileName", params.FileName,
				"error", err.Error(),
			)

			// Report the file as not found.
			return &frontend.FrontendStaticNotFound{}
		}

		return &frontend.FrontendStaticOK{
			Payload: file,
		}
	}
}
