package client

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"

	ksmconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"ivpn.net/auth/services/token/config"
	"ivpn.net/auth/services/token/model"
)

var (
	ErrEmptyInput = "input string cannot be empty"
)

type Signer struct {
	Cfg    *config.Config
	Client *kms.Client
}

func NewSigner(cfg config.Config) (*Signer, error) {
	ctx := context.Background()
	ksmCfg, err := ksmconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &Signer{
		Cfg:    &cfg,
		Client: kms.NewFromConfig(ksmCfg),
	}, nil
}

func (s *Signer) Token(input string) (*model.HSMToken, error) {
	if input == "" {
		return nil, fmt.Errorf("%s", ErrEmptyInput)
	}

	digest := sha512.Sum512([]byte(input))

	if s.Cfg.Mock {
		return &model.HSMToken{
			Token: base64.StdEncoding.EncodeToString(digest[:]),
		}, nil
	}

	signInput := &kms.SignInput{
		KeyId:            &s.Cfg.KeyId,
		Message:          digest[:],
		MessageType:      types.MessageTypeDigest,
		SigningAlgorithm: types.SigningAlgorithmSpecRsassaPssSha256,
	}

	signOut, err := s.Client.Sign(context.Background(), signInput)
	if err != nil {
		return nil, fmt.Errorf("failed to sign input: %w", err)
	}

	return &model.HSMToken{
		Token: base64.StdEncoding.EncodeToString(signOut.Signature),
	}, nil
}
