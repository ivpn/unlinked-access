package config

import "os"

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
}

type Config struct {
	TokenServer TokenServerConfig
	DB          DBConfig
	Service     ServiceConfig
}

func New() (Config, error) {
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
		},
	}, nil
}
