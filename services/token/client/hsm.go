package client

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"ivpn.net/auth/services/token/model"
)

var (
	ErrEmptyInput = "input string cannot be empty"
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

func (h *HSM) Token(input string) (*model.HSMToken, error) {
	if input == "" {
		return nil, fmt.Errorf("%s", ErrEmptyInput)
	}

	keyID := "kms-key-id"
	digest := sha512.Sum512([]byte(input))
	ctx := context.Background()

	signInput := &kms.SignInput{
		KeyId:            &keyID,
		Message:          digest[:],
		MessageType:      kmstypes.MessageTypeDigest,
		SigningAlgorithm: kmstypes.SigningAlgorithmSpecRsassaPssSha256,
	}

	signOut, err := h.Client.Sign(ctx, signInput)
	if err != nil {
		return nil, fmt.Errorf("failed to sign input: %w", err)
	}

	return &model.HSMToken{
		Token: base64.StdEncoding.EncodeToString(signOut.Signature),
	}, nil
}
