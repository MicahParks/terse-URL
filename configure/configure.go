package configure

import (
	"context"

	"github.com/MicahParks/ctxerrgroup"
	"github.com/teris-io/shortid"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/storage"
)

func Configure() (frontendDir string, logger *zap.SugaredLogger, invalidPaths []string, keycloakInfo *KeycloakInfo, shortID *shortid.Shortid, terseStore storage.TerseStore, visitsStore storage.VisitsStore, err error) {

	// Create a logger.
	var zapLogger *zap.Logger
	if zapLogger, err = zap.NewDevelopment(); err != nil { // TODO Make NewProduction
		return "", nil, nil, nil, nil, nil, nil, err
	}
	logger = zapLogger.Sugar()
	logger.Info("Logger created. Starting configuration.")

	// Read the configuration from the environment.
	var config *Configuration
	if config, err = readEnvVars(); err != nil {
		logger.Fatalw("Failed to read configuration from environment variables.",
			"error", err.Error(),
		)
		return "", nil, nil, nil, nil, nil, nil, err // Should be unreachable.
	}

	// Create the Keycloak information data structure.
	keycloakInfo = &KeycloakInfo{
		BaseURL:      config.KeycloakBaseURL,
		ClientID:     config.KeycloakID,
		ClientSecret: config.KeycloakSecret,
		Realm:        config.KeycloakRealm,
	}

	// Set the database timeout.
	defaultTimeout = config.DefaultTimeout

	// Create a channel to report errors asynchronously.
	errChan := make(chan error)

	// Log any errors printed asynchronously.
	go handleAsyncError(errChan, logger)

	// Create a ctxerrgroup for misc asynchronous items needed for requests.
	// TODO Make sure this is used properly.
	group := ctxerrgroup.New(config.WorkerCount, config.WorkersBuffer, true, func(_ ctxerrgroup.Group, err error) {
		logger.Errorw("An error occurred with a ctxerrgroup worker.",
			"error", err.Error(),
		)
	})

	// Create the correct VisitsStore.
	switch config.VisitsStoreType {

	// Use an in memory implementation for Visits storage.
	case memoryStorage:
		visitsStore = storage.NewMemVisits()

	// Use MongoDB for Visits storage.
	case mongoStorage:

		// Create a context that will fail the creation of the Visits store if it takes too long to contact MongoDB.
		ctx, cancel := DefaultCtx()
		defer cancel()

		// Create the Visits storage with MongoDB.
		opts := options.Client().ApplyURI(config.VisitsMongoURI)
		if visitsStore, err = storage.NewMongoDBVisits(ctx, config.VisitsMongoDatabase, config.VisitsMongoCollection, opts); err != nil {
			logger.Fatalw("Failed to reach MongoDB.",
				"store", "VisitsStore",
				"error", err.Error(),
			)
			return "", nil, nil, nil, nil, nil, nil, err // Should be unreachable.
		}

	// If no known Visits storage was specified, don't store visits.
	default:
		visitsStore = nil // Ineffectual assignment, but more clearly shows visits will not be tracked by default.
	}

	// Create the correct TerseStore.
	switch config.TerseStoreType {

	// Use an in memory implementation for Terse storage.
	case memoryStorage:
		storage.NewMemTerse(DefaultCtx, errChan, &group, visitsStore)

	// Set up MongoDB for the Terse storage.
	case mongoStorage:

		// Create a context that will fail the creation of the Terse store if it takes too long to contact MongoDB.
		ctx, cancel := DefaultCtx()
		defer cancel()

		// Create the Terse storage with MongoDB.
		opts := options.Client().ApplyURI(config.TerseMongoURI)
		if terseStore, err = storage.NewMongoDBTerse(ctx, DefaultCtx, config.TerseMongoDatabase, config.TerseMongoCollection, errChan, &group, visitsStore, opts); err != nil {
			logger.Fatalw("Failed to reach MongoDB.",
				"store", "VisitsStore",
				"error", err.Error(),
			)
			return "", nil, nil, nil, nil, nil, nil, err // Should be unreachable.
		}

	// If no known Terse storage was specified in the configuration use an in memory implementation.
	// TODO Change to bbolt when implemented.
	default:
		terseStore = storage.NewMemTerse(DefaultCtx, errChan, &group, visitsStore)
	}

	// Schedule any existing deletions for the Terse pairs.
	ctx, cancel := DefaultCtx()
	defer cancel()
	if err = terseStore.ScheduleDeletions(ctx); err != nil {
		logger.Fatalw("Failed to schedule deletions",
			"error", err.Error(),
		)
		return "", nil, nil, nil, nil, nil, nil, err
	}

	return frontendDir, logger, invalidPaths, keycloakInfo, shortID, terseStore, visitsStore, nil
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
