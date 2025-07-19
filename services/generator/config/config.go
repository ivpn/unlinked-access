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

type Config struct {
	TokenServer TokenServerConfig
	DB          DBConfig
}

func New() (Config, error) {
	return Config{
		TokenServer: TokenServerConfig{
			Host: os.Getenv("TOKEN_HOST"),
			Port: os.Getenv("TOKEN_PORT"),
		},
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}, nil
}
