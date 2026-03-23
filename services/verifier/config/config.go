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
	SampleData         bool
	Mock               bool
	AWSKeyId           string
	AWSAccessKeyId     string
	AWSSecretAccessKey string
	AWSRegion          string
	FortanixEndpoint   string
	FortanixApiKey     string
	FortanixKeyId      string
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
			SampleData:         os.Getenv("SAMPLE_DATA") == "true",
			Mock:               os.Getenv("TOKEN_MOCK") == "true",
			AWSKeyId:           os.Getenv("AWS_TOKEN_KEY_ID"),
			AWSAccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
			AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			AWSRegion:          os.Getenv("AWS_REGION"),
			FortanixEndpoint:   os.Getenv("FORTANIX_ENDPOINT"),
			FortanixApiKey:     os.Getenv("FORTANIX_API_KEY"),
			FortanixKeyId:      os.Getenv("FORTANIX_KEY_ID"),
		},
	}, nil
}
