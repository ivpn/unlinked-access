package client

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"

	ksmconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"ivpn.net/auth/services/verifier/config"
)

type Verifier struct {
	Cfg    *config.Config
	Client *kms.Client
}

func NewVerifierAWS(cfg config.Config) (*Verifier, error) {
	ctx := context.Background()
	kmsCreds := credentials.NewStaticCredentialsProvider(
		cfg.Service.AWSAccessKeyId,
		cfg.Service.AWSSecretAccessKey,
		"",
	)
	ksmCfg, err := ksmconfig.LoadDefaultConfig(
		ctx,
		ksmconfig.WithRegion(cfg.Service.AWSRegion),
		ksmconfig.WithCredentialsProvider(kmsCreds),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &Verifier{
		Cfg:    &cfg,
		Client: kms.NewFromConfig(ksmCfg),
	}, nil
}

func (s *Verifier) Verify(signature string, data []byte) error {
	digest := sha256.Sum256(data)
	digestBase64 := base64.StdEncoding.EncodeToString(digest[:])

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		log.Printf("error decoding signature: %v", err)
		return fmt.Errorf("error decoding signature: %w", err)
	}

	message := sha512.Sum512([]byte(digestBase64))

	verifyInput := &kms.VerifyMacInput{
		KeyId:        &s.Cfg.Service.KeyId,
		Message:      message[:],
		Mac:          sigBytes,
		MacAlgorithm: types.MacAlgorithmSpecHmacSha256,
	}

	verifyOut, _ := s.Client.VerifyMac(context.Background(), verifyInput)
	if verifyOut == nil {
		log.Printf("error verifying manifest signature: verifyOut is nil")
		return fmt.Errorf("error verifying manifest signature: verifyOut is nil")
	}
	if !verifyOut.MacValid {
		log.Printf("manifest signature is invalid")
		return fmt.Errorf("manifest signature is invalid")
	}

	log.Println("manifest signature OK")

	return nil
}
