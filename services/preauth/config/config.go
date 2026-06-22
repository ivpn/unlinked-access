package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

type APIConfig struct {
	AddPort           string
	AddPSK            string
	GetPort           string
	GetPSK            string
	PreauthTTL        time.Duration
	SessionServices   []string
	SessionURLs       []string
	SessionPSKs       []string
	ApiTrustedProxies []string
	ApiAllowIPs       []string
}

type RedisConfig struct {
	Addr             string
	Addrs            []string
	MasterName       string
	Username         string
	Password         string
	FailoverUsername string
	FailoverPassword string
	TLSEnabled       bool
	CertFile         string
	KeyFile          string
	CACertFile       string
}

type TokenServerConfig struct {
	Host          string
	Port          string
	TLSEnabled    bool
	TLSCACertFile string
	TLSCertFile   string
	TLSKeyFile    string
}

type Config struct {
	API         APIConfig
	Redis       RedisConfig
	TokenServer TokenServerConfig
}

func New() (Config, error) {
	preauthTTLStr := os.Getenv("PREAUTH_TTL")
	preauthTTL, err := time.ParseDuration(preauthTTLStr)
	if err != nil {
		return Config{}, err
	}

	sessionServices := strings.Split(os.Getenv("SESSION_SERVICE"), ",")
	sessionURLs := strings.Split(os.Getenv("SESSION_URL"), ",")
	sessionPSKs := strings.Split(os.Getenv("SESSION_PSK"), ",")
	redisAddrs := strings.Split(os.Getenv("REDIS_ADDRESSES"), ",")
	apiTrustedProxies := strings.Split(os.Getenv("API_TRUSTED_PROXIES"), ",")
	apiAllowIPs := strings.Split(os.Getenv("API_ALLOW_IPS"), ",")

	return Config{
		API: APIConfig{
			AddPort:           os.Getenv("PREAUTH_ADD_PORT"),
			AddPSK:            os.Getenv("PREAUTH_ADD_PSK"),
			GetPort:           os.Getenv("PREAUTH_GET_PORT"),
			GetPSK:            os.Getenv("PREAUTH_GET_PSK"),
			PreauthTTL:        preauthTTL,
			SessionServices:   sessionServices,
			SessionURLs:       sessionURLs,
			SessionPSKs:       sessionPSKs,
			ApiTrustedProxies: apiTrustedProxies,
			ApiAllowIPs:       apiAllowIPs,
		},
		Redis: RedisConfig{
			Addr:             os.Getenv("REDIS_ADDR"),
			Addrs:            redisAddrs,
			MasterName:       os.Getenv("REDIS_MASTER_NAME"),
			Username:         os.Getenv("REDIS_USERNAME"),
			Password:         os.Getenv("REDIS_PASSWORD"),
			FailoverUsername: os.Getenv("REDIS_FAILOVER_USERNAME"),
			FailoverPassword: os.Getenv("REDIS_FAILOVER_PASSWORD"),
			TLSEnabled:       os.Getenv("REDIS_TLS_ENABLED") == "true",
			CertFile:         os.Getenv("REDIS_CERT_FILE"),
			KeyFile:          os.Getenv("REDIS_KEY_FILE"),
			CACertFile:       os.Getenv("REDIS_CA_CERT_FILE"),
		},
		TokenServer: TokenServerConfig{
			Host:          os.Getenv("TOKEN_HOST"),
			Port:          os.Getenv("TOKEN_PORT"),
			TLSEnabled:    os.Getenv("TOKEN_TLS_ENABLED") == "true",
			TLSCACertFile: os.Getenv("TOKEN_TLS_CLIENT_CA_FILE"),
			TLSCertFile:   os.Getenv("TOKEN_TLS_CLIENT_CERT_FILE"),
			TLSKeyFile:    os.Getenv("TOKEN_TLS_CLIENT_KEY_FILE"),
		},
	}, nil
}

// Validate checks that all required configuration values are present.
func (c Config) Validate() error {
	required := map[string]string{
		"PREAUTH_ADD_PORT": c.API.AddPort,
		"PREAUTH_ADD_PSK":  c.API.AddPSK,
		"PREAUTH_GET_PORT": c.API.GetPort,
		"PREAUTH_GET_PSK":  c.API.GetPSK,
		"TOKEN_HOST":       c.TokenServer.Host,
		"TOKEN_PORT":       c.TokenServer.Port,
	}
	for name, val := range required {
		if val == "" {
			return errors.New("required env var not set: " + name)
		}
	}
	if c.Redis.Addr == "" && (len(c.Redis.Addrs) == 0 || c.Redis.Addrs[0] == "") {
		return errors.New("required env var not set: REDIS_ADDR or REDIS_ADDRESSES")
	}
	if c.API.PreauthTTL <= 0 {
		return errors.New("PREAUTH_TTL must be a positive duration")
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
