package main

import (
	"errors"
	"testing"

	"ivpn.net/auth/services/token/model"
)

// MockHSMClient is a mock implementation of the HSMClient interface for testing
type MockHSMClient struct {
	mockToken  *model.HSMToken
	mockError  error
	input      string
	ttlMinutes int
}

// Token implements the HSMClient interface for the mock
func (m *MockHSMClient) Token(input string, ttlMinutes int) (*model.HSMToken, error) {
	// Store the parameters for verification
	m.input = input
	m.ttlMinutes = ttlMinutes
	return m.mockToken, m.mockError
}

func TestGenerateToken_Success(t *testing.T) {
	// Arrange
	expectedToken := &model.HSMToken{
		Token: "test-token",
		// Add other fields as needed based on your HSMToken struct
	}
	mockHSM := &MockHSMClient{
		mockToken: expectedToken,
		mockError: nil,
	}

	svc := New(mockHSM)
	inputStr := "test-input"
	ttl := 60

	// Act
	token, err := svc.GenerateToken(inputStr, ttl)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if token != expectedToken {
		t.Errorf("Expected token to be %v, got %v", expectedToken, token)
	}

	if mockHSM.input != inputStr {
		t.Errorf("Expected input to be %v, got %v", inputStr, mockHSM.input)
	}

	if mockHSM.ttlMinutes != ttl {
		t.Errorf("Expected ttlMinutes to be %v, got %v", ttl, mockHSM.ttlMinutes)
	}
}

func TestGenerateToken_Error(t *testing.T) {
	// Arrange
	expectedError := errors.New("hsm error")
	mockHSM := &MockHSMClient{
		mockToken: nil,
		mockError: expectedError,
	}

	svc := New(mockHSM)
	inputStr := "test-input"
	ttl := 30

	// Act
	token, err := svc.GenerateToken(inputStr, ttl)

	// Assert
	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	if token != nil {
		t.Errorf("Expected token to be nil, got %v", token)
	}

	if mockHSM.input != inputStr {
		t.Errorf("Expected input to be %v, got %v", inputStr, mockHSM.input)
	}

	if mockHSM.ttlMinutes != ttl {
		t.Errorf("Expected ttlMinutes to be %v, got %v", ttl, mockHSM.ttlMinutes)
	}
}

func TestGenerateToken_DifferentParameters(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		ttlMinutes int
	}{
		{"Empty input", "", 60},
		{"Zero TTL", "test-input", 0},
		{"Negative TTL", "test-input", -10},
		{"Long input", "this-is-a-very-long-input-string-for-testing", 120},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			expectedToken := &model.HSMToken{Token: "mock-token"}
			mockHSM := &MockHSMClient{
				mockToken: expectedToken,
				mockError: nil,
			}

			svc := New(mockHSM)

			// Act
			_, err := svc.GenerateToken(tc.input, tc.ttlMinutes)

			// Assert
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if mockHSM.input != tc.input {
				t.Errorf("Expected input to be %v, got %v", tc.input, mockHSM.input)
			}

			if mockHSM.ttlMinutes != tc.ttlMinutes {
				t.Errorf("Expected ttlMinutes to be %v, got %v", tc.ttlMinutes, mockHSM.ttlMinutes)
			}
		})
	}
}
