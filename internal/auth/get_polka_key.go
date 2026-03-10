package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetPolkaKey(headers http.Header) (string, error) {
	authorizationHeader := headers.Get("Authorization")

	if authorizationHeader == "" {
		return "", errors.New("No authorization header included")
	}

	if !strings.HasPrefix(authorizationHeader, "ApiKey ") {
		return "", errors.New("Malformed authorization header")
	}

	apiKey := strings.TrimPrefix(authorizationHeader, "ApiKey ")

	return strings.TrimSpace(apiKey), nil

}
