package config

import "os"

type Config struct {
	Host               string
	Port               string
	KeyId              string
	Mock               bool
	AWSAccessKeyId     string
	AWSSecretAccessKey string
	AWSRegion          string
}

func New() (Config, error) {
	return Config{
		Host:               os.Getenv("TOKEN_HOST"),
		Port:               os.Getenv("TOKEN_PORT"),
		KeyId:              os.Getenv("TOKEN_KEY_ID"),
		Mock:               os.Getenv("TOKEN_MOCK") == "true",
		AWSAccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:          os.Getenv("AWS_REGION"),
	}, nil
}
