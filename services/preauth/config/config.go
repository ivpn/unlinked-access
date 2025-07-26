package config

import (
	"os"
	"strings"
	"time"
)

type APIConfig struct {
	Port        string
	PSK         string
	AllowOrigin string
	AllowedIPs  []string
	PreauthTTL  time.Duration
}

type RedisConfig struct {
	Addr                  string
	Addrs                 []string
	MasterName            string
	Username              string
	Password              string
	FailoverUsername      string
	FailoverPassword      string
	TLSEnabled            bool
	CertFile              string
	KeyFile               string
	CACertFile            string
	TLSInsecureSkipVerify bool // Optional: Only for testing, use false in production
}

type TokenServerConfig struct {
	Host string
	Port string
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

	allowedIPs := strings.Split(os.Getenv("PREAUTH_ALLOWED_IPS"), ",")
	redisAddrs := strings.Split(os.Getenv("REDIS_ADDRESSES"), ",")

	return Config{
		API: APIConfig{
			Port:        os.Getenv("PREAUTH_PORT"),
			PSK:         os.Getenv("PREAUTH_PSK"),
			AllowOrigin: os.Getenv("PREAUTH_ALLOW_ORIGIN"),
			AllowedIPs:  allowedIPs,
			PreauthTTL:  preauthTTL,
		},
		Redis: RedisConfig{
			Addr:                  os.Getenv("REDIS_ADDR"),
			Addrs:                 redisAddrs,
			MasterName:            os.Getenv("REDIS_MASTER_NAME"),
			Username:              os.Getenv("REDIS_USERNAME"),
			Password:              os.Getenv("REDIS_PASSWORD"),
			FailoverUsername:      os.Getenv("REDIS_FAILOVER_USERNAME"),
			FailoverPassword:      os.Getenv("REDIS_FAILOVER_PASSWORD"),
			TLSEnabled:            os.Getenv("REDIS_TLS_ENABLED") == "true",
			CertFile:              os.Getenv("REDIS_CERT_FILE"),
			KeyFile:               os.Getenv("REDIS_KEY_FILE"),
			CACertFile:            os.Getenv("REDIS_CA_CERT_FILE"),
			TLSInsecureSkipVerify: os.Getenv("REDIS_TLS_INSECURE_SKIP_VERIFY") == "true",
		},
		TokenServer: TokenServerConfig{
			Host: os.Getenv("TOKEN_HOST"),
			Port: os.Getenv("TOKEN_PORT"),
		},
	}, nil
}
