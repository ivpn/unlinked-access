package client

import (
	"crypto/sha512"
	"fmt"
	"time"

	"ivpn.net/auth/services/token/model"
)

var (
	ErrEmptyInput    = "input string cannot be empty"
	ErrGenerateToken = "failed to generate token"
)

type HSM struct{}

func NewHSM() *HSM {
	return &HSM{}
}

func (h *HSM) Token(input string, ttlMinutes int) (*model.HSMToken, error) {
	// TODO: Implement HSM signing
	if input == "" {
		return nil, fmt.Errorf("%s", ErrEmptyInput)
	}

	inputHash := sha512.Sum512([]byte(input))
	expiresAt := time.Now().Add(time.Duration(ttlMinutes) * time.Minute)

	return &model.HSMToken{
		Token:     string(inputHash[:]),
		ExpiresAt: expiresAt,
	}, nil
}
