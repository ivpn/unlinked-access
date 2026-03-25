package config

import (
	"os"
	"strings"
)

type APIConfig struct {
	Port              string
	PSK               string
	ApiTrustedProxies []string
	ApiAllowIPs       []string
}

type Config struct {
	API APIConfig
}

func New() (Config, error) {
	apiTrustedProxies := strings.Split(os.Getenv("API_TRUSTED_PROXIES"), ",")
	apiAllowIPs := strings.Split(os.Getenv("API_ALLOW_IPS"), ",")

	return Config{
		API: APIConfig{
			Port:              os.Getenv("DISTRIBUTOR_PORT"),
			PSK:               os.Getenv("DISTRIBUTOR_PSK"),
			ApiTrustedProxies: apiTrustedProxies,
			ApiAllowIPs:       apiAllowIPs,
		},
	}, nil
}
