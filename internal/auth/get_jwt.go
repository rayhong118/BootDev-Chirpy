package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("No authorization header included")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("Malformed authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	return strings.TrimSpace(token), nil
}
