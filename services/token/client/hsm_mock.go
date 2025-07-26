package client

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"time"

	"ivpn.net/auth/services/token/model"
)

var (
	ErrEmptyInput    = "input string cannot be empty"
	ErrGenerateToken = "failed to generate token"
)

type MockHSM struct{}

func NewMockHSM() *MockHSM {
	return &MockHSM{}
}

func (h *MockHSM) Token(input string, ttlMinutes int) (*model.HSMToken, error) {
	if input == "" {
		return nil, fmt.Errorf("%s", ErrEmptyInput)
	}

	// Generate a mock HSM secret key (in real HSM, this is securely stored and never exposed)
	secretKey := make([]byte, 32)
	_, err := rand.Read(secretKey)
	if err != nil {
		return nil, fmt.Errorf(ErrGenerateToken+": %v", err)
	}

	// Create an HMAC-SHA512 signature using the secret key and input
	mac := hmac.New(sha512.New, secretKey)
	mac.Write([]byte(input))
	signature := mac.Sum(nil)

	// Encode the signature as a base64 token
	token := base64.URLEncoding.EncodeToString(signature)

	// Set expiration
	expiresAt := time.Now().Add(time.Duration(ttlMinutes) * time.Minute)

	return &model.HSMToken{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}
