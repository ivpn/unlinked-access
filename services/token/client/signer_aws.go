package client

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"

	ksmconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"ivpn.net/auth/services/token/config"
	"ivpn.net/auth/services/token/model"
)

const ErrEmptyInput = "input string cannot be empty"

type SignerAWS struct {
	Cfg    *config.Config
	Client *kms.Client
}

func NewSignerAWS(cfg config.Config) (*SignerAWS, error) {
	ctx := context.Background()
	kmsCreds := credentials.NewStaticCredentialsProvider(
		cfg.AWSAccessKeyId,
		cfg.AWSSecretAccessKey,
		"",
	)
	ksmCfg, err := ksmconfig.LoadDefaultConfig(
		ctx,
		ksmconfig.WithRegion(cfg.AWSRegion),
		ksmconfig.WithCredentialsProvider(kmsCreds),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &SignerAWS{
		Cfg:    &cfg,
		Client: kms.NewFromConfig(ksmCfg),
	}, nil
}

func (s *SignerAWS) Generate(ctx context.Context, input string) (*model.HSMToken, error) {
	if input == "" {
		return nil, fmt.Errorf("%s", ErrEmptyInput)
	}

	digest := sha512.Sum512([]byte(input))

	if s.Cfg.Mock {
		return &model.HSMToken{
			Token: base64.StdEncoding.EncodeToString(digest[:]),
		}, nil
	}

	generateInput := &kms.GenerateMacInput{
		KeyId:        &s.Cfg.AWSKeyId,
		Message:      digest[:],
		MacAlgorithm: types.MacAlgorithmSpecHmacSha256,
	}

	signOut, err := s.Client.GenerateMac(ctx, generateInput)
	if err != nil {
		return nil, fmt.Errorf("failed to sign input: %w", err)
	}

	return &model.HSMToken{
		Token: base64.StdEncoding.EncodeToString(signOut.Mac),
	}, nil
}

func (s *SignerAWS) Authenticate() error {
	return nil
}
