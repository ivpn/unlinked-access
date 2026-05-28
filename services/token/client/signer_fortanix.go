package client

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/fortanix/sdkms-client-go/sdkms"
	"ivpn.net/auth/services/token/config"
	"ivpn.net/auth/services/token/model"
)

type SignerFortanix struct {
	Cfg    *config.Config
	Client *sdkms.Client
}

func NewSignerFortanix(cfg config.Config) (*SignerFortanix, error) {
	httpClient := &http.Client{Timeout: 30 * time.Second}
	client := sdkms.Client{
		Endpoint:   cfg.FortanixEndpoint,
		HTTPClient: httpClient,
	}

	_, err := client.AuthenticateWithAPIKey(context.Background(), cfg.FortanixApiKey)
	if err != nil {
		return nil, err
	}

	return &SignerFortanix{
		Cfg:    &cfg,
		Client: &client,
	}, nil
}

func (s *SignerFortanix) Generate(ctx context.Context, input string) (*model.HSMToken, error) {
	if input == "" {
		return nil, fmt.Errorf("%s", ErrEmptyInput)
	}

	digest := sha512.Sum512([]byte(input))
	data := sdkms.Blob(digest[:])
	keyId := s.Cfg.FortanixKeyId

	if s.Cfg.Mock {
		return &model.HSMToken{
			Token: base64.StdEncoding.EncodeToString(digest[:]),
		}, nil
	}

	alg := sdkms.DigestAlgorithmSha256
	req := sdkms.MacRequest{
		Data: data,
		Alg:  &alg,
		Key:  sdkms.SobjectByID(keyId),
	}

	res, err := s.Client.Mac(ctx, req)
	if err != nil {
		return nil, err
	}

	return &model.HSMToken{
		Token: base64.StdEncoding.EncodeToString(res.Mac),
	}, nil
}

func (s *SignerFortanix) Authenticate() error {
	_, err := s.Client.AuthenticateWithAPIKey(context.Background(), s.Cfg.FortanixApiKey)
	return err
}

func (s *SignerFortanix) Verify(ctx context.Context, data [64]byte, signature string) (bool, error) {
	sigData, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	mac := sdkms.Blob(sigData)
	alg := sdkms.DigestAlgorithmSha256
	keyId := s.Cfg.FortanixKeyId
	req := sdkms.VerifyMacRequest{
		Data: data[:],
		Mac:  &mac,
		Alg:  &alg,
		Key:  sdkms.SobjectByID(keyId),
	}

	res, err := s.Client.MacVerify(ctx, req)
	if err != nil {
		return false, err
	}

	return res.Result, nil
}
