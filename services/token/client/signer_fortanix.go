package client

import (
	"context"
	"net/http"

	"github.com/fortanix/sdkms-client-go/sdkms"
	"ivpn.net/auth/services/token/config"
)

type FortanixSigner struct {
	Cfg    *config.Config
	Client *sdkms.Client
}

func NewFortanixSigner(cfg config.Config) (*FortanixSigner, error) {
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
