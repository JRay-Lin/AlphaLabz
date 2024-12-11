package routes

import (
	"errors"
	"net/http"
	"strings"
)

func ExtractAuthToken(r *http.Request) (string, error) {
	// Extract the token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	// Check if the header starts with "Bearer "
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("invalid authorization format")
	}

	// Extract the token
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}
