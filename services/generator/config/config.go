package config

import (
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
	Host string
	Port string
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
			Host: os.Getenv("TOKEN_HOST"),
			Port: os.Getenv("TOKEN_PORT"),
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
