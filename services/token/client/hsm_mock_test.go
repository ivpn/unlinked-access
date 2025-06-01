package client

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"ivpn.net/auth/services/token/model"
)

func TestMockHSMClient_Token(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		ttl       int
		wantErr   bool
		errorText string
	}{
		{
			name:      "Valid input",
			input:     "test-input",
			ttl:       60,
			wantErr:   false,
			errorText: "",
		},
		{
			name:      "Empty input",
			input:     "",
			ttl:       60,
			wantErr:   true,
			errorText: ErrEmptyInput,
		},
		{
			name:      "Zero TTL",
			input:     "test-input",
			ttl:       0,
			wantErr:   false,
			errorText: "",
		},
		{
			name:      "Negative TTL",
			input:     "test-input",
			ttl:       -10,
			wantErr:   false,
			errorText: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New()

			// Note the time before generating the token
			beforeTime := time.Now()

			got, err := h.Token(tt.input, tt.ttl)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errorText) {
					t.Errorf("Expected error containing %q but got %q", tt.errorText, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify token is not empty
			if got.Token == "" {
				t.Error("Expected non-empty token")
			}

			// Verify token is properly base64 URL encoded
			if !isBase64URLEncoded(got.Token) {
				t.Errorf("Token is not properly base64 URL encoded: %s", got.Token)
			}

			// Verify token expiration time
			expectedExpiry := beforeTime.Add(time.Duration(tt.ttl) * time.Minute)
			tolerance := 2 * time.Second // Allow a small time difference

			if got.ExpiresAt.Sub(expectedExpiry) > tolerance || expectedExpiry.Sub(got.ExpiresAt) > tolerance {
				t.Errorf("Expected expiration around %v but got %v", expectedExpiry, got.ExpiresAt)
			}

			// Test that two different inputs generate different tokens
			if tt.input != "" {
				anotherToken, _ := h.Token(tt.input+"different", tt.ttl)
				if got.Token == anotherToken.Token {
					t.Error("Different inputs should generate different tokens")
				}
			}
		})
	}
}

func TestMockHSMClient_TokenConsistency(t *testing.T) {
	// Test that same input generates different tokens due to random secret key
	h := New()
	token1, err := h.Token("same-input", 60)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	token2, err := h.Token("same-input", 60)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if token1.Token == token2.Token {
		t.Error("Expected different tokens for same input due to random secret key")
	}
}

func TestMockHSMClient_ReturnedType(t *testing.T) {
	h := New()
	result, err := h.Token("test", 30)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify the returned type is what we expect
	if _, ok := interface{}(result).(*model.HSMToken); !ok {
		t.Error("Expected return type to be *model.HSMToken")
	}
}

// Helper function to check if a string is base64 URL encoded
func isBase64URLEncoded(s string) bool {
	// Check if string only contains valid base64 URL characters
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') || c == '-' || c == '_' || c == '=') {
			return false
		}
	}

	// Try to decode
	_, err := base64.URLEncoding.DecodeString(s)
	return err == nil
}
