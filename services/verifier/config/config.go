package config

import (
	"os"
)

type APIConfig struct {
	ManifestURL string
	ManifestPSK string
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type Config struct {
	API APIConfig
	DB  DBConfig
}

func New() (Config, error) {
	return Config{
		API: APIConfig{
			ManifestURL: os.Getenv("MANIFEST_URL"),
			ManifestPSK: os.Getenv("MANIFEST_PSK"),
		},
		DB: DBConfig{
			Host:     os.Getenv("CLIENT_DB_HOST"),
			Port:     os.Getenv("CLIENT_DB_PORT"),
			Name:     os.Getenv("CLIENT_DB_NAME"),
			User:     os.Getenv("CLIENT_DB_USER"),
			Password: os.Getenv("CLIENT_DB_PASSWORD"),
		},
	}, nil
}
