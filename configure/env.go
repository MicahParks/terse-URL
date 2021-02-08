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

	// defaultPrefix is the default HTTP prefix for all shortened URLs.
	defaultPrefix = "https://terseurl.com/"

	// defaultTemplatePath is the default path in the file system to look for the HTML template that will be used for
	// JavaScript and HTML meta tag redirects, along with Javascript fingerprinting.
	defaultTemplatePath = "redirect.gohtml"

	// defaultWorkerCount is the default amount of workers to have in the ctxerrgroup.
	defaultWorkerCount = 4
)

var (

	// ErrCantBeZeroOrNegative indicates that an integer value cannot be zero negative, but a zero or negative value was
	// provided.
	ErrCantBeZeroOrNegative = errors.New("integer cannot be negative")

	// alwaysInvalidPaths
	alwaysInvalidPaths = []string{"api", "docs", "frontend", "favicon.ico", "swagger.json", "robots.txt"}

	// defaultTimeout is the default timeout for any incoming (from clients) and outgoing (to databases) requests.
	defaultTimeout = time.Minute
)

// configuration holds all the necessary information for
type configuration struct {
	DefaultTimeout   time.Duration
	InvalidPaths     []string
	Prefix           string
	ShortIDParanoid  bool
	ShortIDSeed      uint64
	TemplatePath     string
	SummaryStoreJSON string
	TerseStoreJSON   string
	VisitsStoreJSON  string
	WorkerCount      uint
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
	config.Prefix = os.Getenv("HTTP_PREFIX")
	if config.Prefix == "" {
		config.Prefix = defaultPrefix
	}
	config.ShortIDParanoid = len(os.Getenv("SHORTID_PARANOID")) != 0
	config.TemplatePath = os.Getenv("TEMPLATE_PATH")
	config.TerseStoreJSON = os.Getenv("TERSE_STORE_JSON")
	config.VisitsStoreJSON = os.Getenv("VISITS_STORE_JSON")
	config.SummaryStoreJSON = os.Getenv("SUMMARY_STORE_JSON")

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
