package service

import (
	"errors"
	"testing"

	"ivpn.net/auth/services/token/config"
	"ivpn.net/auth/services/token/model"
)

// MockHSMClient is a mock implementation of the HSMClient interface for testing
type MockHSMClient struct {
	mockToken *model.HSMToken
	mockError error
	input     string
}

// Token implements the HSMClient interface for the mock
func (m *MockHSMClient) Generate(input string) (*model.HSMToken, error) {
	// Store the parameters for verification
	m.input = input
	return m.mockToken, m.mockError
}

func (m *MockHSMClient) Authenticate() error {
	return nil
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

	cfg, err := config.New()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	svc := New(mockHSM, cfg)
	inputStr := "test-input"

	// Act
	token, err := svc.generateToken(inputStr)

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
}

func TestGenerateToken_Error(t *testing.T) {
	// Arrange
	expectedError := errors.New("hsm error")
	mockHSM := &MockHSMClient{
		mockToken: nil,
		mockError: expectedError,
	}

	cfg, err := config.New()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	svc := New(mockHSM, cfg)
	inputStr := "test-input"

	// Act
	token, err := svc.generateToken(inputStr)

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
}

func TestGenerateToken_DifferentParameters(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Empty input", ""},
		{"Zero TTL", "test-input"},
		{"Negative TTL", "test-input"},
		{"Long input", "this-is-a-very-long-input-string-for-testing"},
	}

	cfg, err := config.New()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			expectedToken := &model.HSMToken{Token: "mock-token"}
			mockHSM := &MockHSMClient{
				mockToken: expectedToken,
				mockError: nil,
			}

			svc := New(mockHSM, cfg)

			// Act
			_, err := svc.generateToken(tc.input)

			// Assert
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if mockHSM.input != tc.input {
				t.Errorf("Expected input to be %v, got %v", tc.input, mockHSM.input)
			}
		})
	}
}
