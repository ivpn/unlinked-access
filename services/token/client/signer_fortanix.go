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

type SignerFortanix struct {
	Cfg    *config.Config
	Client *sdkms.Client
}

func NewSignerFortanix(cfg config.Config) (*SignerFortanix, error) {
	client := sdkms.Client{
		Endpoint:   cfg.FortanixEndpoint,
		HTTPClient: http.DefaultClient,
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

func (s *SignerFortanix) Generate(input string) (*model.HSMToken, error) {
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

func (s *SignerFortanix) Verify(data [64]byte, signature string) (bool, error) {
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

	res, err := s.Client.MacVerify(context.Background(), req)
	if err != nil {
		return false, err
	}

	return res.Result, nil
}
