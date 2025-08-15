package config

import "os"

type Config struct {
	Host  string
	Port  string
	KeyId string
	Mock  bool
}

func New() (Config, error) {
	return Config{
		Host:  os.Getenv("TOKEN_HOST"),
		Port:  os.Getenv("TOKEN_PORT"),
		KeyId: os.Getenv("TOKEN_KEY_ID"),
		Mock:  os.Getenv("TOKEN_MOCK") == "true",
	}, nil
}
