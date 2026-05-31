package config

import (
	"errors"
	"os"
	"strconv"
)

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type TokenServerConfig struct {
	Host          string
	Port          string
	TLSEnabled    bool
	TLSCACertFile string
	TLSCertFile   string
	TLSKeyFile    string
}

type ServiceConfig struct {
	SampleData bool
	Mock       bool
	TPS        int // total TPS across all goroutines
}

type Config struct {
	TokenServer TokenServerConfig
	DB          DBConfig
	Service     ServiceConfig
}

func New() (Config, error) {
	tps, err := strconv.Atoi(os.Getenv("GENERATOR_TPS"))
	if err != nil {
		return Config{}, err
	}

	return Config{
		TokenServer: TokenServerConfig{
			Host:          os.Getenv("TOKEN_HOST"),
			Port:          os.Getenv("TOKEN_PORT"),
			TLSEnabled:    os.Getenv("TOKEN_TLS_ENABLED") == "true",
			TLSCACertFile: os.Getenv("TOKEN_TLS_CLIENT_CA_FILE"),
			TLSCertFile:   os.Getenv("TOKEN_TLS_CLIENT_CERT_FILE"),
			TLSKeyFile:    os.Getenv("TOKEN_TLS_CLIENT_KEY_FILE"),
		},
		DB: DBConfig{
			Host:     os.Getenv("SERVER_DB_HOST"),
			Port:     os.Getenv("SERVER_DB_PORT"),
			Name:     os.Getenv("SERVER_DB_NAME"),
			User:     os.Getenv("SERVER_DB_USER"),
			Password: os.Getenv("SERVER_DB_PASSWORD"),
		},
		Service: ServiceConfig{
			SampleData: os.Getenv("SAMPLE_DATA") == "true",
			Mock:       os.Getenv("GENERATOR_MOCK") == "true",
			TPS:        tps,
		},
	}, nil
}

// Validate checks that all required configuration values are present.
func (c Config) Validate() error {
	required := map[string]string{
		"TOKEN_HOST":         c.TokenServer.Host,
		"TOKEN_PORT":         c.TokenServer.Port,
		"SERVER_DB_HOST":     c.DB.Host,
		"SERVER_DB_PORT":     c.DB.Port,
		"SERVER_DB_NAME":     c.DB.Name,
		"SERVER_DB_USER":     c.DB.User,
		"SERVER_DB_PASSWORD": c.DB.Password,
	}
	for name, val := range required {
		if val == "" {
			return errors.New("required env var not set: " + name)
		}
	}
	if c.Service.TPS <= 0 {
		return errors.New("GENERATOR_TPS must be a positive integer")
	}
	if c.TokenServer.TLSEnabled {
		if c.TokenServer.TLSCACertFile == "" {
			return errors.New("required env var not set: TOKEN_TLS_CLIENT_CA_FILE")
		}
		if c.TokenServer.TLSCertFile == "" {
			return errors.New("required env var not set: TOKEN_TLS_CLIENT_CERT_FILE")
		}
		if c.TokenServer.TLSKeyFile == "" {
			return errors.New("required env var not set: TOKEN_TLS_CLIENT_KEY_FILE")
		}
	}
	return nil
}
