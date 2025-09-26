package client

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	ksmconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
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

func NewAWSSigner(cfg config.Config) (*Signer, error) {
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

	return &Signer{
		Cfg:    &cfg,
		Client: kms.NewFromConfig(ksmCfg),
	}, nil
}

func (s *Signer) Token(input string) (*model.HSMToken, error) {
	start := time.Now()

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
		KeyId:        &s.Cfg.KeyId,
		Message:      digest[:],
		MacAlgorithm: types.MacAlgorithmSpecHmacSha256,
	}

	signOut, err := s.Client.GenerateMac(context.Background(), generateInput)
	if err != nil {
		return nil, fmt.Errorf("failed to sign input: %w", err)
	}

	elapsed := time.Since(start)
	log.Printf("Token() completed in %s", elapsed)

	return &model.HSMToken{
		Token: base64.StdEncoding.EncodeToString(signOut.Mac),
	}, nil
}
