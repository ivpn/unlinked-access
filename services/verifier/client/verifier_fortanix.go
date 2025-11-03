package client

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/fortanix/sdkms-client-go/sdkms"
	"ivpn.net/auth/services/verifier/config"
)

type VerifierFortanix struct {
	Cfg    *config.Config
	Client *sdkms.Client
}

func NewVerifierFortanix(cfg config.Config) (*VerifierFortanix, error) {
	client := sdkms.Client{
		Endpoint:   cfg.Service.FortanixEndpoint,
		HTTPClient: http.DefaultClient,
	}

	_, err := client.AuthenticateWithAPIKey(context.Background(), cfg.Service.FortanixApiKey)
	if err != nil {
		return nil, err
	}

	return &VerifierFortanix{
		Cfg:    &cfg,
		Client: &client,
	}, nil
}

func (s *VerifierFortanix) Verify(signature string, data []byte) error {
	digest := sha256.Sum256(data)
	digestBase64 := base64.StdEncoding.EncodeToString(digest[:])

	if s.Cfg.Service.Mock {
		hash512 := sha512.Sum512([]byte(digestBase64))
		digestBase64 = base64.StdEncoding.EncodeToString(hash512[:])

		if digestBase64 != signature {
			return fmt.Errorf("invalid manifest signature (mock)")
		}

		log.Println("manifest signature (mock) OK")

		return nil
	}

	sigData, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("error decoding signature: %w", err)
	}

	message := sha512.Sum512([]byte(digestBase64))
	mac := sdkms.Blob(sigData)
	alg := sdkms.DigestAlgorithmSha256
	keyId := s.Cfg.Service.FortanixKeyId
	req := sdkms.VerifyMacRequest{
		Data: message[:],
		Mac:  &mac,
		Alg:  &alg,
		Key:  sdkms.SobjectByID(keyId),
	}

	res, err := s.Client.MacVerify(context.Background(), req)
	if err != nil {
		return err
	}

	if !res.Result {
		return fmt.Errorf("invalid manifest signature")
	}

	log.Println("manifest signature OK")

	return nil
}
