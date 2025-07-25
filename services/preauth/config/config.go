package config

import (
	"os"
	"strings"
)

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
	Redis       RedisConfig
	TokenServer TokenServerConfig
}

func New() (Config, error) {
	redisAddrs := strings.Split(os.Getenv("REDIS_ADDRESSES"), ",")

	return Config{
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
