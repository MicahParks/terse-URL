package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/MicahParks/jwks"
	"github.com/dgrijalva/jwt-go"

	"github.com/MicahParks/terseurl/models"
)

var (

	// ErrClaims indicates the JWT is did not have the expected claims.
	ErrClaims = errors.New("the JWT did not have the expected claims")

	// ErrInvalidJWT indicates the JWT is invalid.
	ErrInvalidJWT = errors.New("the JWT is invalid")
)

// JWTHandler is a function signature that takes in a Base64 encoded JWT and returns the auth principal from it.
type JWTHandler func(jwtB64 string) (principal *models.Principal, err error)

// HandleJWT creates a JWT auth handler via a closure.
//
// TODO Add logging. Error is returned to user. Log error. Generic thing back to user.
func HandleJWT(ks jwks.Keystore) (authHandler JWTHandler) {
	return func(jwtB64 string) (principal *models.Principal, err error) {

		// Remove the "Bearer " prefix.
		jwtB64 = strings.TrimPrefix(jwtB64, "Bearer ")

		// Parse the JWT.
		token, err := jwt.Parse(jwtB64, ks.KeyFunc())
		if err != nil {
			return nil, fmt.Errorf("failed to parse token: %w", err)
		}

		// Confirm the JWT is valid.
		if !token.Valid {
			return nil, ErrInvalidJWT
		}

		// Get the claims from the JWT.
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("couldn't type assert JWT map claims: %w", ErrClaims)
		}

		// Transform the claim back to JSON.
		var data []byte
		if data, err = json.Marshal(claims); err != nil {
			return nil, fmt.Errorf("failed to marshal claims back to JSON type: %w: %v", ErrClaims, err)
		}

		// Unmarshal the claims JSON into the principal.
		principal = &models.Principal{}
		if err = json.Unmarshal(data, principal); err != nil {
			return nil, fmt.Errorf("failed to unmarshal claims back from JSON: %w: %v", ErrClaims, err)
		}

		return principal, nil
	}
}
