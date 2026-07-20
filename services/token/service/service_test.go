package service

import (
	"context"
	"errors"
	"testing"

	proto "ivpn.net/auth/services/proto"
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
func (m *MockHSMClient) GenerateToken(ctx context.Context, input string) (*model.HSMToken, error) {
	// Store the parameters for verification
	m.input = input
	return m.mockToken, m.mockError
}

func (m *MockHSMClient) GenerateSignature(ctx context.Context, input string) (*model.HSMToken, error) {
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
	inputStr := "i-TEST-1234-ABCD"

	// Act
	token, err := svc.GenerateToken(context.Background(), &proto.Request{Input: inputStr})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if token.GetToken() != expectedToken.Token {
		t.Errorf("Expected token to be %v, got %v", expectedToken.Token, token.GetToken())
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
	inputStr := "i-TEST-1234-ABCD"

	// Act
	token, err := svc.GenerateToken(context.Background(), &proto.Request{Input: inputStr})

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

func TestGenerateToken_InvalidFormat(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Empty input", ""},
		{"No prefix", "ABCD-1234-EFGH"},
		{"Lowercase", "i-abcd-1234-efgh"},
		{"Too short", "i-ABC-1234-EFGH"},
		{"Too long", "i-ABCDE-1234-EFGH"},
	}

	cfg, err := config.New()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockHSM := &MockHSMClient{}
			svc := New(mockHSM, cfg)

			_, err := svc.GenerateToken(context.Background(), &proto.Request{Input: tc.input})
			if err == nil {
				t.Errorf("Expected error for input %q, got nil", tc.input)
			}
		})
	}
}

func TestGenerateToken_DifferentParameters(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Alphanumeric", "i-ABCD-1234-EFGH"},
		{"All digits", "i-1234-5678-9012"},
		{"All letters", "i-ABCD-EFGH-IJKL"},
		{"Mixed", "i-A1B2-C3D4-E5F6"},
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
			_, err := svc.GenerateToken(context.Background(), &proto.Request{Input: tc.input})

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
