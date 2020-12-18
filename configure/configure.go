package configure

import (
	"os"

	"go.uber.org/zap"

	"github.com/MicahParks/shakesearch"
)

const (

	// defaultWorksPath is the default path to find the file that contains the complete works of Shakespeare in text.
	defaultWorksPath = "completeworks.txt"

	// worksPathEnvVar represent the environment variable used to find the path to the file that contains complete works
	// of Shakespeare in text. It will use the default value if none is given.
	worksPathEnvVar = "SHAKESPEARES_WORKS"
)

// Configure gathers all required information and creates the required Go structs to run the service.
func Configure() (logger *zap.SugaredLogger, shakeSearcher *shakesearch.ShakeSearcher, err error) {

	// Create a logger.
	var zapLogger *zap.Logger
	if zapLogger, err = zap.NewDevelopment(); err != nil { // TODO Make NewProduction.
		return nil, nil, err
	}
	logger = zapLogger.Sugar()
	logger.Info("Logger created. Starting configuration.")

	// Get the complete works of Shakespeare's file path from an environment variable.
	worksPath := os.Getenv(worksPathEnvVar)
	if worksPath == "" {
		worksPath = defaultWorksPath
	}

	// Create the ShakeSearcher.
	if shakeSearcher, err = shakesearch.NewShakeSearcher(worksPath); err != nil {
		logger.Fatalw("Failed to create the ShakeSearcher.",
			"error", err.Error(),
		)
		return nil, nil, err // Should be unreachable.
	}

	return logger, shakeSearcher, nil
}
