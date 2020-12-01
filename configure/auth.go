package configure

import (
	"errors"
	"strings"

	"go.uber.org/zap"

	"github.com/Nerzal/gocloak/v7"

	"github.com/MicahParks/terse-URL/models"
)

const (

	// headerPrefix is the prefix the JWT has in the "Authorization" header's value.
	headerPrefix = "Bearer "
)

var (

	// ErrJWTExpired indicates the JWT has expired.
	ErrJWTExpired = errors.New("JWT has expired")
)

type KeycloakInfo struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	Realm        string
}

func HandleAuth(keycloakInfo *KeycloakInfo, logger *zap.SugaredLogger) (authHandler func(headerValue string) (jwtInfo *models.JWTInfo, err error), err error) {

	// Create the Keycloak client.
	keycloak := gocloak.NewClient(keycloakInfo.BaseURL)

	// Create a context for logging into Keycloak.
	ctx, cancel := DefaultCtx()
	defer cancel()

	// Log into Keycloak with the given client.
	if _, err = keycloak.LoginClient(ctx, keycloakInfo.ClientID, keycloakInfo.ClientSecret, keycloakInfo.Realm); err != nil {
		return nil, err
	}

	// Create the authentication handler via a closure.
	return func(headerValue string) (jwtInfo *models.JWTInfo, err error) {

		// Strip the prefix from the header.
		headerValue = strings.TrimPrefix(headerValue, headerPrefix)

		// Create the JWT info structure that will be passed to the endpoints.
		jwtInfo = &models.JWTInfo{}

		// Create a context to authorize this JWT.
		ctx, cancel := DefaultCtx()
		defer cancel()

		// Check the JWT with Keycloak.
		var res *gocloak.RetrospecTokenResult
		if res, err = keycloak.RetrospectToken(ctx, headerValue, keycloakInfo.ClientID, keycloakInfo.ClientSecret, keycloakInfo.Realm); err != nil {
			logger.Errorw("Failed to retrospect token with Keycloak.",
				"error", err.Error(),
			)
			return nil, err
		}

		// Confirm the JWT is active.
		if !*res.Active {
			logger.Warnw("An inactive JWT was received.",
				"JWT", headerValue,
			)
			return nil, ErrJWTExpired
		}

		// TODO Populate the JWT information in the return value.

		return jwtInfo, nil
	}, nil
}
