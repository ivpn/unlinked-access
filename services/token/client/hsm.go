package client

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"ivpn.net/auth/services/token/model"
)

var (
	ErrEmptyInput    = "input string cannot be empty"
	ErrGenerateToken = "failed to generate token"
)

type HSM struct {
	Client *kms.Client
}

func NewHSM() (*HSM, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &HSM{
		Client: kms.NewFromConfig(cfg),
	}, nil
}

func (h *HSM) Token(input string, ttlMinutes int) (*model.HSMToken, error) {
	// TODO: Implement HSM signing
	if input == "" {
		return nil, fmt.Errorf("%s", ErrEmptyInput)
	}

	inputHash := sha512.Sum512([]byte(input))
	expiresAt := time.Now().Add(time.Duration(ttlMinutes) * time.Minute)

	return &model.HSMToken{
		Token:     base64.StdEncoding.EncodeToString(inputHash[:]),
		ExpiresAt: expiresAt,
	}, nil
}
