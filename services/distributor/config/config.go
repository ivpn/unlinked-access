package config

import (
	"os"
	"strings"
)

type APIConfig struct {
	Port         string
	PSK          string
	AllowOrigins string
	AllowedIPs   []string
}

type Config struct {
	API APIConfig
}

func New() (Config, error) {
	allowedIPs := strings.Split(os.Getenv("DISTRIBUTOR_ALLOWED_IPS"), ",")

	return Config{
		API: APIConfig{
			Port:         os.Getenv("DISTRIBUTOR_PORT"),
			PSK:          os.Getenv("DISTRIBUTOR_PSK"),
			AllowOrigins: os.Getenv("DISTRIBUTOR_ALLOW_ORIGINS"),
			AllowedIPs:   allowedIPs,
		},
	}, nil
}
