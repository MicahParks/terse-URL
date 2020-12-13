package configure

import (
	"context"

	"github.com/MicahParks/ctxerrgroup"
	"github.com/teris-io/shortid"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/storage"
)

// Configuration is the Go structure that contains all needed configurations gathered on startup.
type Configuration struct {
	Logger          *zap.SugaredLogger
	InvalidPaths    []string
	ShortID         *shortid.Shortid
	ShortIDParanoid bool
	TerseStore      storage.TerseStore
	VisitsStore     storage.VisitsStore
}

// Configure gathers all startup configurations, formats them, and returns them as a Go struct.
func Configure() (config Configuration, err error) {

	// Create a logger.
	var zapLogger *zap.Logger
	if zapLogger, err = zap.NewDevelopment(); err != nil { // TODO Make NewProduction
		return Configuration{}, err
	}
	logger := zapLogger.Sugar()
	logger.Info("Logger created. Starting configuration.")
	config.Logger = logger

	// Read the configuration from the environment.
	var rawConfig *configuration
	if rawConfig, err = readEnvVars(); err != nil {
		logger.Fatalw("Failed to read configuration from environment variables.",
			"error", err.Error(),
		)
		return Configuration{}, err // Should be unreachable.
	}

	// Set the database timeout.
	defaultTimeout = rawConfig.DefaultTimeout

	// Create a channel to report errors asynchronously.
	errChan := make(chan error)

	// Log any errors printed asynchronously.
	go handleAsyncError(errChan, logger)

	// Create a ctxerrgroup for misc asynchronous items needed for requests.
	group := ctxerrgroup.New(rawConfig.WorkerCount, func(_ ctxerrgroup.Group, err error) {
		logger.Errorw("An error occurred with a ctxerrgroup worker.",
			"error", err.Error(),
		)
	})

	// Create the correct VisitsStore.
	switch rawConfig.VisitsStoreType {

	// Use an in memory implementation for Visits storage.
	case memoryStorage:
		config.VisitsStore = storage.NewMemVisits()

	//// Use MongoDB for Visits storage.
	//case mongoStorage:
	//
	//	// Create a context that will fail the creation of the Visits store if it takes too long to contact MongoDB.
	//	ctx, cancel := DefaultCtx()
	//	defer cancel()
	//
	//	// Create the Visits storage with MongoDB.
	//	opts := options.Client().ApplyURI(rawConfig.VisitsMongoURI)
	//	if config.VisitsStore, err = storage.NewMongoDBVisits(ctx, rawConfig.VisitsMongoDatabase, rawConfig.VisitsMongoCollection, opts); err != nil {
	//		logger.Fatalw("Failed to reach MongoDB.",
	//			"store", "VisitsStore",
	//			"error", err.Error(),
	//		)
	//		return Configuration{}, err // Should be unreachable.
	//	}

	// If no known Visits storage was specified, don't store visits.
	default:
		config.VisitsStore = nil // Ineffectual assignment, but more clearly shows visits will not be tracked by default.
	}

	// Create the correct TerseStore.
	switch rawConfig.TerseStoreType {

	// Use an in memory implementation for Terse storage.
	case memoryStorage:
		config.TerseStore = storage.NewMemTerse(DefaultCtx, errChan, &group, config.VisitsStore)

	//// Set up MongoDB for the Terse storage.
	//case mongoStorage:
	//
	//	// Create a context that will fail the creation of the Terse store if it takes too long to contact MongoDB.
	//	ctx, cancel := DefaultCtx()
	//	defer cancel()
	//
	//	// Create the Terse storage with MongoDB.
	//	opts := options.Client().ApplyURI(rawConfig.TerseMongoURI)
	//	if config.TerseStore, err = storage.NewMongoDBTerse(ctx, DefaultCtx, rawConfig.TerseMongoDatabase, rawConfig.TerseMongoCollection, errChan, &group, config.VisitsStore, opts); err != nil {
	//		logger.Fatalw("Failed to reach MongoDB.",
	//			"store", "VisitsStore",
	//			"error", err.Error(),
	//		)
	//		return Configuration{}, err // Should be unreachable.
	//	}

	// If no known Terse storage was specified in the configuration use an in memory implementation.
	default:
		config.TerseStore = storage.NewMemTerse(DefaultCtx, errChan, &group, config.VisitsStore)
	}

	// Create the short ID generator.
	if config.ShortID, err = shortid.New(1, shortid.DefaultABC, rawConfig.ShortIDSeed); err != nil {
		return Configuration{}, err
	}

	// Copy over any other needed raw config info.
	config.InvalidPaths = rawConfig.InvalidPaths
	config.ShortIDParanoid = rawConfig.ShortIDParanoid

	return config, nil
}

// DefaultCtx creates a context and its cancel function using the default timeout or one provided during configuration.
func DefaultCtx() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}

// handleAsyncError logs errors asynchronously for an error channel.
func handleAsyncError(errChan <-chan error, logger *zap.SugaredLogger) {
	var err error
	for err = range errChan {
		logger.Errorw("An error has been reported asynchronously.",
			"error", err.Error(),
		)
	}
}
