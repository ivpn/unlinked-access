package client

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/fortanix/sdkms-client-go/sdkms"
	"ivpn.net/auth/services/token/config"
	"ivpn.net/auth/services/token/model"
)

type FortanixSigner struct {
	Cfg    *config.Config
	Client *sdkms.Client
}

func NewSignerFortanix(cfg config.Config) (*FortanixSigner, error) {
	client := sdkms.Client{
		Endpoint:   cfg.FortanixEndpoint,
		HTTPClient: http.DefaultClient,
	}

	_, err := client.AuthenticateWithAPIKey(context.Background(), cfg.FortanixApiKey)
	if err != nil {
		return nil, err
	}

	return &FortanixSigner{
		Cfg:    &cfg,
		Client: &client,
	}, nil
}

func (s *FortanixSigner) Token(input string) (*model.HSMToken, error) {
	if input == "" {
		return nil, fmt.Errorf("%s", ErrEmptyInput)
	}

	// start := time.Now()

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

	res, err := s.Client.Mac(context.Background(), req)
	if err != nil {
		return nil, err
	}

	// elapsed := time.Since(start)
	// log.Printf("Token() completed in %s", elapsed)

	return &model.HSMToken{
		Token: base64.StdEncoding.EncodeToString(res.Mac),
	}, nil
}
