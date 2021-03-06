// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/MicahParks/terseurl/auth"
	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/endpoints"
	"github.com/MicahParks/terseurl/endpoints/public"
	"github.com/MicahParks/terseurl/endpoints/system"
	"github.com/MicahParks/terseurl/middleware"
	"github.com/MicahParks/terseurl/models"
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

	// Check to see if auth is turned on.
	if config.UseAuth {

		// Create a context for creating the JWT handler.
		ctx, cancel := context.WithTimeout(context.TODO(), time.Minute*5) // TODO Make configurable.
		sleep := time.Second * 5                                          // TODO Make configurable.

		// Configure the JWT auth.
		//
		// TODO Make HTTP client configurable?
		api.JWTAuth, err = auth.HandleJWT(ctx, nil, config.JWKSURL, logger.Named("JWT Authenticator"), sleep)
		if err != nil {
			logger.Fatalw("failed to get JWKS", // TODO Remove.
				"error", err.Error(),
			)
		}
		cancel()
		logger.Info("Authentication with JWKS configured.")
	} else {
		api.JWTAuth = func(s string) (*models.Principal, error) {
			return nil, nil
		}
		logger.Info("Authentication is turned off.")
	}

	// Assign the endpoint handlers.
	api.APIExportHandler = endpoints.HandleExport(logger.Named("POST /api/export"), config.StoreManager)
	api.APIFrontendMetaHandler = endpoints.HandleMeta(logger.Named("POST /api/frontend/meta"))
	api.APIImportHandler = endpoints.HandleImport(logger.Named("POST /api/import"), config.StoreManager)
	api.APIShortenedDeleteHandler = endpoints.HandleShortenedDelete(logger.Named("DELETE /api/shortened"), config.StoreManager)
	api.APIShortenedPrefixHandler = endpoints.HandleShortenedPrefix(logger.Named("POST /api/prefix"), config.Prefix)
	api.APIShortenedSummaryHandler = endpoints.HandleShortenedSummary(logger.Named("POST /api/summary"), config.StoreManager)
	api.APITerseReadHandler = endpoints.HandleTerseRead(logger.Named("POST /api/terse"), config.StoreManager)
	api.APITerseWriteHandler = endpoints.HandleWrite(logger.Named("POST /api/write/{operation}"), config.ShortID, config.StoreManager)
	api.APIVisitsDeleteHandler = endpoints.HandlerVisitsDelete(logger.Named("DELETE /api/visits"), config.StoreManager)
	api.APIVisitsReadHandler = endpoints.HandleVisitsRead(logger.Named("POST /api/visits"), config.StoreManager)
	api.PublicPublicRedirectHandler = public.HandleRedirect(logger.Named("GET /{shortenedURL}"), config.Template, config.StoreManager)
	api.SystemSystemAliveHandler = system.HandleAlive()

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {

		// Close the error channel to clean up the async error logging goroutine.
		defer close(config.ErrChan)

		// Create a context to close the Terse and Visits stores.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Close the TerseStore.
		if err = config.StoreManager.Close(ctx); err != nil {
			logger.Errorw("Failed to close the TerseStore.",
				"error", err.Error(),
			)
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
	limit.SetIPLookups([]string{"X-Forwarded-For"}) // TODO Add string for regular lookup.

	// Set up the rate limiter middleware.
	toll := tollbooth.LimitHandler(limit, handler) // TODO Logging middleware. Maybe another rate limiter instead.

	// Set up the frontend middleware.
	frontendMiddleware, err := middleware.FrontendMiddleware(toll)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Follow the HTTP middleware pattern.
	return frontendMiddleware
}
