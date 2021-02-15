// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/endpoints"
	"github.com/MicahParks/terseurl/endpoints/frontend"
	"github.com/MicahParks/terseurl/endpoints/public"
	"github.com/MicahParks/terseurl/endpoints/system"
	"github.com/MicahParks/terseurl/restapi/operations"
)

//go:generate swagger generate server --target ../../terseurl --name TerseURL --spec ../swagger.yml --principal interface{}

func configureFlags(api *operations.TerseurlAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.TerseurlAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Configure the service.
	config, err := configure.Configure()
	if err != nil {
		log.Fatalf("Failed to configure the service.\nError: %s\n", err.Error())
	}
	logger := config.Logger

	api.UseSwaggerUI()

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	// Create the HTML producer.
	api.HTMLProducer = configure.HTMLProducer(logger)

	// Assign the endpoint handlers.
	api.APIFrontendMetaHandler = frontend.HandleMeta(logger.Named("/api/frontend/meta"))
	api.APITerseDeleteHandler = endpoints.HandleDelete(logger.Named("/api/delete"), config.TerseStore)
	api.APITerseDeleteSomeHandler = endpoints.HandleDeleteSome(logger.Named("/api/delete/some"), config.TerseStore)
	api.APITerseExportHandler = endpoints.HandleExport(logger.Named("/api/export"), config.TerseStore)
	api.APITerseExportSomeHandler = endpoints.HandleExportSome(logger.Named("/api/export/some"), config.TerseStore)
	api.APITerseImportHandler = endpoints.HandleImport(logger.Named("/api/import"), config.TerseStore)
	api.APITersePrefixHandler = endpoints.HandlePrefix(logger.Named("/api/prefix"), config.Prefix)
	api.APITerseSummaryHandler = endpoints.HandleSummary(logger.Named("/api/summary"), config.SummaryStore)
	api.APITerseTerseHandler = endpoints.HandleTerse(logger.Named("/api/terse/{shortened}"), config.TerseStore)
	api.APITerseVisitsHandler = endpoints.HandleVisits(logger.Named("/api/visits/{shortened}"), config.VisitsStore)
	api.APITerseWriteHandler = endpoints.HandleWrite(logger.Named("/api/write/{operation}"), config.ShortID, config.TerseStore)
	api.PublicTerseRedirectHandler = public.HandleRedirect(logger.Named("/{shortened}"), config.Template, config.TerseStore)
	api.SystemAliveHandler = system.HandleAlive()

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {

		// Close the error channel to clean up the async error logging goroutine.
		defer close(config.ErrChan)

		// Create a context to close the Terse and Visits stores.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Close the TerseStore.
		if err = config.TerseStore.Close(ctx); err != nil {
			logger.Errorw("Failed to close the TerseStore.",
				"error", err.Error(),
			)
		}

		// Close the VisitsStore.
		if config.VisitsStore != nil {
			if err = config.VisitsStore.Close(ctx); err != nil {
				logger.Errorw("Failed to close the VisitsStore.",
					"error", err.Error(),
				)
			}
		}
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {

	// Create an incoming request rate limiter that only allows 1 request per section and forgets about clients after 1
	// hour.
	limit := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})

	// Find the IP of the client in the X-Forwarded-For header, because Caddy will be the server in front of this.
	limit.SetIPLookups([]string{"X-Forwarded-For"})

	// Follow the HTTP middleware pattern.
	return tollbooth.LimitHandler(limit, handler) // TODO Logging middleware. Maybe another rate limiter instead.
}