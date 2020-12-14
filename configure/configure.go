package configure

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/MicahParks/ctxerrgroup"
	"github.com/teris-io/shortid"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/storage"
)

const (

	// configPathTerseStore is the location to find the TerseStore JSON configuration file.
	configPathTerseStore = "terseStore.json"

	// configPathVisitsStore is the location to find the VisitsStore JSON configuration file.
	configPathVisitsStore = "visitsStore.json"
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

	// Get the VisitsStore configuration.
	var visitsConfig json.RawMessage
	if visitsConfig, err = readStorageConfig(rawConfig.VisitsStoreJSON, logger, configPathVisitsStore); err != nil {
		return Configuration{}, err
	}

	// Create the VisitsStore.
	var visitsStoreType string
	if config.VisitsStore, visitsStoreType, err = storage.NewVisitsStore(visitsConfig); err != nil {
		logger.Fatalw("Failed to create VisitsStore.",
			"type", visitsStoreType,
			"error", err.Error(),
		)
		return Configuration{}, err // Should be unreachable.
	}
	logger.Infow("Created VisitsStore.",
		"type", visitsStoreType,
	)

	// Get the TerseStore configuration.
	var terseConfig json.RawMessage
	if terseConfig, err = readStorageConfig(rawConfig.TerseStoreJSON, logger, configPathTerseStore); err != nil {
		return Configuration{}, err
	}

	// Create the TerseStore.
	var terseStoreType string
	if config.TerseStore, terseStoreType, err = storage.NewTerseStore(terseConfig, DefaultCtx, errChan, &group, config.VisitsStore); err != nil {
		logger.Fatalw("Failed to create TerseStore.",
			"type", terseStoreType,
			"error", err.Error(),
		)
		return Configuration{}, err // Should be unreachable.
	}
	logger.Infow("Created TerseStore.",
		"type", terseStoreType,
	)

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

// TODO
func readStorageConfig(envValue string, logger *zap.SugaredLogger, configPath string) (configJSON json.RawMessage, err error) {

	// Decide if the configPath is valid. Generate a long message from it.
	var logMessage string
	switch configPath {
	case configPathTerseStore:
		logMessage = "TerseStore"
	case configPathVisitsStore:
		logMessage = "VisitsStore"
	default:
		panic("not implemented")
	}

	// Use the environment variable's value, if present.
	if envValue != "" {

		// Log that no environment variable was present.
		logger.Infow(fmt.Sprintf("No %s environment variable configuration present. Attempting to read configuration file.", logMessage),
			"filePath", configPath,
		)

		// Read the JSON file where the configuration is expected to be at.
		var data []byte
		if data, err = ioutil.ReadFile(configPath); err != nil {
			return nil, err
		}

		// Place the data in the envValue variable.
		envValue = string(data)
	} else {

		// Log that the environment variable is present.
		logger.Infow(fmt.Sprintf("%s environment variable configuration present.", logMessage))
	}

	return json.RawMessage(envValue), nil
}
