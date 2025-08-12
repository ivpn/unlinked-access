package client

import (
	"strings"
	"testing"

	"ivpn.net/auth/services/token/model"
)

func TestMockHSM_Token(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		errorText string
	}{
		{
			name:      "Valid input",
			input:     "test-input",
			wantErr:   false,
			errorText: "",
		},
		{
			name:      "Empty input",
			input:     "",
			wantErr:   true,
			errorText: ErrEmptyInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, _ := NewMockHSM()

			got, err := h.Token(tt.input)

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

			// Test that two different inputs generate different tokens
			if tt.input != "" {
				anotherToken, _ := h.Token(tt.input + "different")
				if got.Token == anotherToken.Token {
					t.Error("Different inputs should generate different tokens")
				}
			}
		})
	}
}

func TestMockHSM_TokenConsistency(t *testing.T) {
	h, _ := NewMockHSM()
	token1, err := h.Token("same-input")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	token2, err := h.Token("same-input")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if token1.Token != token2.Token {
		t.Error("Expected same tokens for same input")
	}
}

func TestMockHSM_ReturnedType(t *testing.T) {
	h, _ := NewMockHSM()
	result, err := h.Token("test")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify the returned type is what we expect
	if _, ok := any(result).(*model.HSMToken); !ok {
		t.Error("Expected return type to be *model.HSMToken")
	}
}
