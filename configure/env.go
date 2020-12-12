package configure

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (

	// defaultMongoDatabase is the default MongoDB database to use.
	defaultMongoDatabase = "terseURL"

	// defaultMongoTerseCollection is the default MongoDB collection name for the TerseStore.
	defaultMongoTerseCollection = "terseStore"

	// defaultMongoVisitsCollection is the default MongoDB collection name for the VisitsStore.
	defaultMongoVisitsCollection = "visitsStore"

	// defaultWorkerCount is the default amount of workers to have in the ctxerrgroup.
	defaultWorkerCount = 4

	// memoryStorage is the constant used when describing a storage backend only in memory.
	memoryStorage = "memory"

	// mongoStorage is the string constant used when describing a storage backend of MongoDB.
	mongoStorage = "mongo"
)

var (

	// ErrCantBeZeroOrNegative indicates that an integer value cannot be zero negative, but a zero or negative value was
	// provided.
	ErrCantBeZeroOrNegative = errors.New("integer cannot be negative")

	// ErrMissingRequiredConfig indicates that a required configuration was missing.
	ErrMissingRequiredConfig = errors.New("required configuration missing")

	// alwaysInvalidPaths
	alwaysInvalidPaths = []string{"api", "docs", "frontend", "swagger.json"}

	// defaultTimeout is the default timeout for any incoming (from clients) and outgoing (to databases) requests.
	defaultTimeout = time.Minute
)

// configuration holds all the necessary information for
type configuration struct {
	DefaultTimeout        time.Duration
	KeycloakBaseURL       string
	KeycloakID            string
	KeycloakRealm         string
	KeycloakSecret        string
	InvalidPaths          []string
	ShortIDSeed           uint64
	TerseMongoCollection  string
	TerseMongoDatabase    string
	TerseMongoURI         string
	TerseStoreType        string
	VisitsMongoCollection string
	VisitsMongoDatabase   string
	VisitsMongoURI        string
	VisitsStoreType       string
	WorkerCount           uint
}

// invalidPathsParse parses a comma separated string into a slice of strings. It adds in paths that are always invalid
// to the slice if not given.
func invalidPathsParse(s string) (invalidPaths []string) {

	// Create the invalid paths slice.
	invalidPaths = make([]string, 0)

	// Iterate through the split string and append it to the slice.
	for _, path := range strings.Split(s, ",") {
		path = strings.TrimSpace(path)
		invalidPaths = append(invalidPaths, path)
	}

	// Make sure all the always invalid paths are in the slice.
	var have bool
	for _, alwaysInvalid := range alwaysInvalidPaths { // TODO Validate.
		have = false
		for _, path := range invalidPaths {
			if alwaysInvalid == path {
				have = true
				break
			}
		}
		if !have {
			invalidPaths = append(invalidPaths, alwaysInvalid)
		}
	}

	return invalidPaths
}

// readEnvVars reads in the environment variables and handles defaults for everything except storage types.
func readEnvVars() (config *configuration, err error) {

	// Create the configuration structure.
	config = &configuration{}

	// Transform the required environment variables to seconds.
	incomingRequestTimeout := os.Getenv("DEFAULT_TIMEOUT")
	if config.DefaultTimeout, err = stringToSeconds(incomingRequestTimeout, defaultTimeout); err != nil {
		return nil, fmt.Errorf("%w: %s", err, incomingRequestTimeout)
	}

	// Transform the required value environment variables into unsigned integers.
	workerCount := os.Getenv("WORKER_COUNT")
	if config.WorkerCount, err = stringToUint(workerCount, defaultWorkerCount); err != nil {
		return nil, fmt.Errorf("%w: %s", err, workerCount)
	}

	// Transform the short ID seed into a uint64, if given.
	shortIDSeed := os.Getenv("SHORTID_SEED")
	if shortIDSeed == "" {
		config.ShortIDSeed = uint64(time.Now().UnixNano())
	} else if config.ShortIDSeed, err = strconv.ParseUint(shortIDSeed, 10, 64); err != nil {
		return nil, fmt.Errorf("could not parse shortid seed: %w", err)
	}

	// Transform the required environment variables to slices.
	config.InvalidPaths = invalidPathsParse(os.Getenv("INVALID_PATHS"))

	// Assign the string value configurations.
	config.KeycloakBaseURL = os.Getenv("KEYCLOAK_BASE_URL")
	config.KeycloakID = os.Getenv("KEYCLOAK_ID")
	config.KeycloakRealm = os.Getenv("KEYCLOAK_REALM")
	config.KeycloakSecret = os.Getenv("KEYCLOAK_SECRET")
	config.TerseMongoCollection = os.Getenv("TERSE_MONGO_COLLECTION")
	config.TerseMongoDatabase = os.Getenv("TERSE_MONGO_DATABASE")
	config.TerseMongoURI = os.Getenv("TERSE_MONGO_URI")
	config.TerseStoreType = os.Getenv("TERSE_STORE_TYPE")
	config.VisitsMongoCollection = os.Getenv("VISITS_MONGO_COLLECTION")
	config.VisitsMongoDatabase = os.Getenv("VISITS_MONGO_DATABASE")
	config.VisitsMongoURI = os.Getenv("VISITS_MONGO_URI")
	config.VisitsStoreType = os.Getenv("VISITS_STORE_TYPE")

	// Confirm none of the Keycloak environment variables are empty.
	if config.KeycloakBaseURL == "" || config.KeycloakID == "" || config.KeycloakRealm == "" || config.KeycloakSecret == "" {
		return nil, fmt.Errorf("%w: All Keycloak enviornment variables must be populated", ErrMissingRequiredConfig)
	}

	// If using MongoDB for Terse storage, check for defaults to use.
	if config.TerseStoreType == mongoStorage {
		if config.TerseMongoCollection == "" {
			config.TerseMongoCollection = defaultMongoTerseCollection
		}
		if config.TerseMongoDatabase == "" {
			config.TerseMongoDatabase = defaultMongoDatabase
		}
		if config.TerseMongoURI == "" {
			return nil, fmt.Errorf("%w: Terse MongoDB URI required when Terse storage is in MongoDB", ErrMissingRequiredConfig)
		}
	}

	// If using MongoDB for visits, check for defaults to use.
	if config.VisitsStoreType == mongoStorage {
		if config.VisitsMongoCollection == "" {
			config.VisitsMongoCollection = defaultMongoVisitsCollection
		}
		if config.VisitsMongoDatabase == "" {
			config.VisitsMongoDatabase = defaultMongoDatabase
		}
		if config.TerseMongoURI == "" {
			return nil, fmt.Errorf("%w: Visits MongoDB URI required when Visits storage is in MongoDB", ErrMissingRequiredConfig)
		}
	}

	return config, nil
}

// stringToSeconds converts a string to a time.Duration. If the string is empty, it uses returns the default.
func stringToSeconds(s string, defaultSeconds time.Duration) (seconds time.Duration, err error) {

	// If not provided, use the default quantity of seconds.
	if s == "" {
		seconds = defaultSeconds
	} else {

		// Convert the string of seconds to the correct Go type.
		var u uint
		if u, err = stringToUint(s, 0); err != nil {
			return 0, err
		}

		// Convert to the correct Go type.
		seconds = time.Second * time.Duration(u)
	}

	return seconds, nil
}

// stringToUint converts a string to an unsigned integer. If the string is empty, it returns the default.
func stringToUint(s string, defaultUint uint) (u uint, err error) {

	// If not provided, use the default unsigned integer.
	if s == "" {
		u = defaultUint
	} else {

		// Convert to an integer.
		var integer int
		if integer, err = strconv.Atoi(s); err != nil {
			return 0, err
		}

		// Confirm the integer is not zero or negative.
		if integer <= 0 {
			return 0, ErrCantBeZeroOrNegative
		}

		// Convert the integer to an unsigned one.
		u = uint(integer)
	}

	return u, nil
}
