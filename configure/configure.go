package configure

import (
	"context"
	"html/template"
	"io/fs"
	"io/ioutil"

	"github.com/MicahParks/ctxerrgroup"
	"github.com/teris-io/shortid"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl"
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
	StaticFS        fs.FS
	StoreManager    storage.StoreManager
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

	// Create the HTML template.
	if config.Template, err = createTemplate(rawConfig.TemplatePath, ""); err != nil {
		return Configuration{}, err
	}

	// Create the file system for the frontend static assets.
	if config.StaticFS, err = terseurl.FrontendFS(rawConfig.StaticFSDirName); err != nil {
		logger.Fatalw("Failed to configure file system for static frontend assets.",
			"error", err.Error(),
		)
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

	// Create the Terse, Visits, and Summary data stores.
	if err = createStores(&config, group, logger, rawConfig); err != nil {
		logger.Fatalw("Failed to create data store.",
			"error", err.Error(),
		)
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

	// Check if the embedded template should be used.
	var tmplStr string
	if filePath == "" {

		// Use the embedded template.
		tmplStr = terseurl.RedirectTemplate
	} else {

		// Read the give template file from the OS.
		var fileData []byte
		if fileData, err = ioutil.ReadFile(filePath); err != nil {
			return nil, err
		}
		tmplStr = string(fileData)
	}

	// Create the Go template from the file.
	return template.New(name).Parse(tmplStr)
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
