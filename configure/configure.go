package configure

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/MicahParks/ctxerrgroup"
	"github.com/teris-io/shortid"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/storage"
)

const (

	// configPathTerseStore is the location to find the TerseStore JSON configuration file.
	configPathTerseStore = "terseStore.json"

	// configPathVisitsStore is the location to find the VisitsStore JSON configuration file.
	configPathVisitsStore = "visitsStore.json"

	// configPathSummaryStore is the location to find the SummaryStore JSON configuration file.
	configPathSummaryStore = "summaryStore.json"
)

// Configuration is the Go structure that contains all needed configurations gathered on startup.
type Configuration struct {
	ErrChan         chan error
	Template        *template.Template
	Logger          *zap.SugaredLogger
	InvalidPaths    []string
	Prefix          string
	ShortID         *shortid.Shortid
	ShortIDParanoid bool
	SummaryStore    storage.SummaryStore
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

	// Figure out the path to the HTML template.
	if rawConfig.TemplatePath == "" {
		rawConfig.TemplatePath = defaultTemplatePath
	}

	// Create the HTML template.
	if config.Template, err = createTemplate(rawConfig.TemplatePath, ""); err != nil {
		return Configuration{}, err
	}

	// Set the database timeout.
	defaultTimeout = rawConfig.DefaultTimeout

	// Create a channel to report errors asynchronously.
	config.ErrChan = make(chan error)

	// Log any errors printed asynchronously.
	go handleAsyncError(config.ErrChan, logger)

	// Create a ctxerrgroup for misc asynchronous items needed for requests.
	group := ctxerrgroup.New(rawConfig.WorkerCount, func(_ ctxerrgroup.Group, err error) {
		logger.Errorw("An error occurred with a ctxerrgroup worker.",
			"error", err.Error(),
		)
	})

	// Get the SummaryStore configuration.
	var summaryConfig json.RawMessage
	if summaryConfig, err = readStorageConfig(rawConfig.SummaryStoreJSON, logger, configPathSummaryStore); err != nil {
		return Configuration{}, err
	}

	// Create the SummaryStore.
	var summaryStoreType string
	if config.SummaryStore, summaryStoreType, err = storage.NewSummaryStore(summaryConfig, config.VisitsStore); err != nil {
		logger.Fatalw("Failed to create SummaryStore.",
			"type", summaryStoreType,
			"error", err.Error(),
		)
		return Configuration{}, err
	}
	logger.Infow("Created SummaryStore.",
		"type", summaryStoreType,
	)

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
	if config.TerseStore, terseStoreType, err = storage.NewTerseStore(terseConfig, DefaultCtx, config.ErrChan, &group, config.SummaryStore, config.VisitsStore); err != nil {
		logger.Fatalw("Failed to create TerseStore.",
			"type", terseStoreType,
			"error", err.Error(),
		)
		return Configuration{}, err // Should be unreachable.
	}
	logger.Infow("Created TerseStore.",
		"type", terseStoreType,
	)

	// TODO ctxCreator
	var summaries map[string]models.TerseSummary
	if summaries, err = storage.InitializeSummaries(ctx, config.TerseStore, config.VisitsStore); err != nil {
		logger.Fatalw("Failed to initialize the summaries for the SummaryStore.",
			"error", err.Error(),
		) // TODO Not fatal.
	}

	// TODO ctxCreator
	if err = config.SummaryStore.Import(ctx, summaries); err != nil {
		logger.Fatalw("Failed to import the initial summaries into the SummaryStore.",
			"error", err.Error(),
		) // TODO Not fatal?
	}

	// Create the short ID generator.
	if config.ShortID, err = shortid.New(1, shortid.DefaultABC, rawConfig.ShortIDSeed); err != nil { // TODO Configure worker count?
		return Configuration{}, err
	}

	// Copy over any other needed raw config info.
	config.InvalidPaths = rawConfig.InvalidPaths
	config.ShortIDParanoid = rawConfig.ShortIDParanoid
	config.Prefix = rawConfig.Prefix

	return config, nil
}

// DefaultCtx creates a context and its cancel function using the default timeout or one provided during configuration.
func DefaultCtx() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}

// createTemplate reads the template file at filePath and turns it into a Golang template with the given name.
func createTemplate(filePath, name string) (tmpl *template.Template, err error) {

	// Read the template file into memory.
	var fileData []byte
	if fileData, err = ioutil.ReadFile(filePath); err != nil {
		return nil, err
	}

	// Create the Go template from the file.
	return template.New(name).Parse(string(fileData))
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

// readStorageConfig determines which value to use. Either the value at the file at configPath or the envValue. The
// chosen value will be turned into a raw JSON message.
func readStorageConfig(envValue string, logger *zap.SugaredLogger, configPath string) (configJSON json.RawMessage, err error) {

	// Decide if the configPath is valid. Generate a long message from it.
	var logMessage string
	switch configPath {
	case configPathSummaryStore:
		logMessage = "SummaryStore"
	case configPathTerseStore:
		logMessage = "TerseStore"
	case configPathVisitsStore:
		logMessage = "VisitsStore"
	default:
		panic("not implemented")
	}

	// Use the environment variable's value, if present.
	if envValue == "" {

		// Log that no environment variable was present.
		message := fmt.Sprintf("No %s environment variable configuration present. Attempting to read configuration file.", logMessage)
		logger.Infow(message,
			"filePath", configPath,
		)

		// Stat the config file.
		var data []byte
		if _, existErr := os.Stat(configPath); existErr != nil {

			// If it doesn't exist. Use the default config.
			if os.IsNotExist(existErr) {

				// Do nothing. Return an empty config.

			} else {

				// Return any other errors.
				return nil, existErr
			}

		} else {
			// Read the JSON file where the configuration is expected to be at.
			if data, err = ioutil.ReadFile(configPath); err != nil {
				return nil, err
			}
		}

		// Place the data in the envValue variable.
		envValue = string(data)
	} else {

		// Log that the environment variable is present.
		logger.Infow(fmt.Sprintf("%s environment variable configuration present.", logMessage))
	}

	return json.RawMessage(envValue), nil
}
