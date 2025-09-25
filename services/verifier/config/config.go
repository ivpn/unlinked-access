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

type ServiceConfig struct {
	KeyId              string
	SampleData         bool
	Mock               bool
	AWSAccessKeyId     string
	AWSSecretAccessKey string
	AWSRegion          string
}

type Config struct {
	API     APIConfig
	DB      DBConfig
	Service ServiceConfig
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
		Service: ServiceConfig{
			KeyId:              os.Getenv("TOKEN_KEY_ID"),
			SampleData:         os.Getenv("SAMPLE_DATA") == "true",
			Mock:               os.Getenv("TOKEN_MOCK") == "true",
			AWSAccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
			AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			AWSRegion:          os.Getenv("AWS_REGION"),
		},
	}, nil
}
