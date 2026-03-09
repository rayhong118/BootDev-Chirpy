package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		headers       http.Header
		expectedToken string
		expectError   bool
	}{
		{
			name: "Valid token",
			headers: http.Header{
				"Authorization": []string{"Bearer valid_token_123"},
			},
			expectedToken: "valid_token_123",
			expectError:   false,
		},
		{
			name: "Valid token with extra spaces",
			headers: http.Header{
				"Authorization": []string{"Bearer    valid_token_123   "},
			},
			expectedToken: "valid_token_123",
			expectError:   false,
		},
		{
			name:          "No authorization header",
			headers:       http.Header{},
			expectedToken: "",
			expectError:   true,
		},
		{
			name: "Malformed header (no Bearer)",
			headers: http.Header{
				"Authorization": []string{"Basic valid_token_123"},
			},
			expectedToken: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)
			if tt.expectError && err == nil {
				t.Errorf("expected an error, but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("did not expect an error, but got: %v", err)
			}
			if token != tt.expectedToken {
				t.Errorf("expected token %q, got %q", tt.expectedToken, token)
			}
		})
	}
}
