package frontend

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/meta"
	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/restapi/operations/api"
)

// HandleMeta creates and /api/frontend/meta endpoint handler via a closure. It will assist the frontend by gathering
// relevant HTML meta information for social media link previews.
func HandleMeta(logger *zap.SugaredLogger) api.FrontendMetaHandlerFunc {
	return func(params api.FrontendMetaParams) middleware.Responder {

		// Debug info.
		logger.Infow("Gathering relevant HTML meta information.",
			"originalURL", params.OriginalURL,
		)

		// Create a request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the HTML meta relevant to social media link previews from the URL.
		og, twitter, err := meta.Get(ctx, params.OriginalURL)
		if err != nil {

			// Log at the appropriate level. Assign the response code and message.
			message := "Failed to perform HTTP GET for social media link preview inheritance."
			logger.Infow(message,
				"originalURL", params.OriginalURL,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := 500
			resp := &api.FrontendMetaDefault{Payload: &models.Error{
				Code:    int64(code),
				Message: message,
			}}
			resp.SetStatusCode(code)
			return resp
		}

		return &api.FrontendMetaOK{
			Payload: &api.FrontendMetaOKBody{
				Og:      og,
				Twitter: twitter,
			},
		}
	}
}
