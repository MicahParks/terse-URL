package configure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/MicahParks/ctxerrgroup"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/storage"
)

// createStores handles the process of creating the SummaryStore, TerseStore, and VisitsStore.
func createStores(config *Configuration, group *ctxerrgroup.Group, logger *zap.SugaredLogger, rawConfig *configuration) (err error) {

	// Get the SummaryStore configuration.
	var summaryConfig json.RawMessage
	if summaryConfig, err = readStorageConfig(rawConfig.SummaryStoreJSON, logger, configPathSummaryStore); err != nil {
		return err
	}

	// Create the SummaryStore.
	var summaryStoreType string
	if config.SummaryStore, summaryStoreType, err = storage.NewSummaryStore(summaryConfig, config.VisitsStore); err != nil {
		logger.Fatalw("Failed to create SummaryStore.",
			"type", summaryStoreType,
			"error", err.Error(),
		)
		return err
	}
	logger.Infow("Created SummaryStore.",
		"type", summaryStoreType,
	)

	// Get the VisitsStore configuration.
	var visitsConfig json.RawMessage
	if visitsConfig, err = readStorageConfig(rawConfig.VisitsStoreJSON, logger, configPathVisitsStore); err != nil {
		return err
	}

	// Create the VisitsStore.
	var visitsStoreType string
	if config.VisitsStore, visitsStoreType, err = storage.NewVisitsStore(visitsConfig); err != nil {
		logger.Fatalw("Failed to create VisitsStore.",
			"type", visitsStoreType,
			"error", err.Error(),
		)
		return err // Should be unreachable.
	}
	logger.Infow("Created VisitsStore.",
		"type", visitsStoreType,
	)

	// Get the TerseStore configuration.
	var terseConfig json.RawMessage
	if terseConfig, err = readStorageConfig(rawConfig.TerseStoreJSON, logger, configPathTerseStore); err != nil {
		return err
	}

	// Create the TerseStore.
	var terseStoreType string
	if config.TerseStore, terseStoreType, err = storage.NewTerseStore(terseConfig, DefaultCtx, config.ErrChan, group, config.SummaryStore, config.VisitsStore); err != nil {
		logger.Fatalw("Failed to create TerseStore.",
			"type", terseStoreType,
			"error", err.Error(),
		)
		return err // Should be unreachable.
	}
	logger.Infow("Created TerseStore.",
		"type", terseStoreType,
	)

	// Read from the Terse store and Visits store to initialize the Summary data.
	ctx, cancel := DefaultCtx()
	var summaries map[string]models.TerseSummary
	if summaries, err = storage.InitializeSummaries(ctx, config.TerseStore, config.VisitsStore); err != nil {
		logger.Fatalw("Failed to initialize the summaries for the SummaryStore.",
			"error", err.Error(),
		)
	}
	cancel()

	// Initialize the Summary data in the Summary data store.
	ctx, cancel = DefaultCtx()
	if err = config.SummaryStore.Import(ctx, summaries); err != nil {
		logger.Fatalw("Failed to import the initial summaries into the SummaryStore.",
			"error", err.Error(),
		)
	}
	cancel()

	return nil
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
