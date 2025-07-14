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
			name:          "valid bearer token",
			headers:       http.Header{"Authorization": []string{"Bearer abc123"}},
			expectedToken: "abc123",
			expectError:   false,
		},
		{
			name:          "valid bearer token with jwt",
			headers:       http.Header{"Authorization": []string{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"}},
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expectError:   false,
		},
		{
			name:          "missing authorization header",
			headers:       http.Header{},
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "empty authorization header",
			headers:       http.Header{"Authorization": []string{""}},
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "bearer token with spaces",
			headers:       http.Header{"Authorization": []string{"Bearer token with spaces"}},
			expectedToken: "token with spaces",
			expectError:   false,
		},
		{
			name:          "bearer token with special characters",
			headers:       http.Header{"Authorization": []string{"Bearer token!@#$%^&*()"}},
			expectedToken: "token!@#$%^&*()",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if token != tt.expectedToken {
				t.Errorf("expected token %q, got %q", tt.expectedToken, token)
			}
		})
	}
}

func TestGetBearerTokenEdgeCases(t *testing.T) {
	t.Run("authorization header without Bearer prefix", func(t *testing.T) {
		headers := http.Header{"Authorization": []string{"Basic abc123"}}

		_, err := GetBearerToken(headers)
		if err == nil {
			t.Error("expected error for non-Bearer authorization header")
		}

		expectedError := "authorization header must start with 'Bearer '"
		if err.Error() != expectedError {
			t.Errorf("expected error %q, got %q", expectedError, err.Error())
		}
	})

	t.Run("just 'Bearer' without token", func(t *testing.T) {
		headers := http.Header{"Authorization": []string{"Bearer"}}

		_, err := GetBearerToken(headers)
		if err == nil {
			t.Error("expected error for Bearer without token")
		}

		expectedError := "authorization header must start with 'Bearer '"
		if err.Error() != expectedError {
			t.Errorf("expected error %q, got %q", expectedError, err.Error())
		}
	})

	t.Run("Bearer with only space", func(t *testing.T) {
		headers := http.Header{"Authorization": []string{"Bearer "}}

		_, err := GetBearerToken(headers)
		if err == nil {
			t.Error("expected error for empty bearer token")
		}

		expectedError := "bearer token is empty"
		if err.Error() != expectedError {
			t.Errorf("expected error %q, got %q", expectedError, err.Error())
		}
	})
}
