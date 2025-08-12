package config

import "os"

type Config struct {
	Host  string
	Port  string
	KeyId string
}

func New() (Config, error) {
	return Config{
		Host:  os.Getenv("TOKEN_HOST"),
		Port:  os.Getenv("TOKEN_PORT"),
		KeyId: os.Getenv("TOKEN_KEY_ID"),
	}, nil
}
