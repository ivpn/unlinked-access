package config

import "os"

type Config struct {
	Host string
	Port string
}

func New() (Config, error) {
	return Config{
		Host: os.Getenv("TOKEN_HOST"),
		Port: os.Getenv("TOKEN_PORT"),
	}, nil
}
