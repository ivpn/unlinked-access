package config

import (
	"errors"
	"os"
)

type Config struct {
	Host               string
	Port               string
	Mock               bool
	AWSKeyId           string
	AWSAccessKeyId     string
	AWSSecretAccessKey string
	AWSRegion          string
	FortanixEndpoint   string
	FortanixApiKey     string
	FortanixKeyId      string
	TLSEnabled         bool
	TLSCertFile        string
	TLSKeyFile         string
	TLSCAFile          string
	Debug              bool
}

func New() (Config, error) {
	return Config{
		Host:               os.Getenv("TOKEN_HOST"),
		Port:               os.Getenv("TOKEN_PORT"),
		Mock:               os.Getenv("TOKEN_MOCK") == "true",
		AWSKeyId:           os.Getenv("AWS_TOKEN_KEY_ID"),
		AWSAccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:          os.Getenv("AWS_REGION"),
		FortanixEndpoint:   os.Getenv("FORTANIX_ENDPOINT"),
		FortanixApiKey:     os.Getenv("FORTANIX_API_KEY"),
		FortanixKeyId:      os.Getenv("FORTANIX_KEY_ID"),
		TLSEnabled:         os.Getenv("TOKEN_TLS_ENABLED") == "true",
		TLSCertFile:        os.Getenv("TOKEN_TLS_CERT_FILE"),
		TLSKeyFile:         os.Getenv("TOKEN_TLS_KEY_FILE"),
		TLSCAFile:          os.Getenv("TOKEN_TLS_CA_FILE"),
		Debug:              os.Getenv("TOKEN_DEBUG") == "true",
	}, nil
}

// Validate returns an error if required credentials are missing for the active signer mode.
func (c *Config) Validate() error {
	if c.Port == "" {
		return errors.New("TOKEN_PORT is required")
	}
	if c.Mock {
		return nil
	}
	// Fortanix is the active signer (NewSignerFortanix is called from main.go).
	// Validate the fields it needs at runtime so the process fails fast at startup.
	if c.FortanixEndpoint == "" {
		return errors.New("FORTANIX_ENDPOINT is required when TOKEN_MOCK=false")
	}
	if c.FortanixApiKey == "" {
		return errors.New("FORTANIX_API_KEY is required when TOKEN_MOCK=false")
	}
	if c.FortanixKeyId == "" {
		return errors.New("FORTANIX_KEY_ID is required when TOKEN_MOCK=false")
	}
	if c.TLSEnabled {
		if c.TLSCertFile == "" {
			return errors.New("TOKEN_TLS_CERT_FILE is required when TOKEN_TLS_ENABLED=true")
		}
		if c.TLSKeyFile == "" {
			return errors.New("TOKEN_TLS_KEY_FILE is required when TOKEN_TLS_ENABLED=true")
		}
		if c.TLSCAFile == "" {
			return errors.New("TOKEN_TLS_CA_FILE is required when TOKEN_TLS_ENABLED=true")
		}
	}
	return nil
}
