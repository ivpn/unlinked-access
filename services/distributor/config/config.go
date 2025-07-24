package config

import "os"

type APIConfig struct {
	Port           string
	PSK            string
	PSKAllowOrigin string
	LogFile        string
}

type Config struct {
	API APIConfig
}

func New() (Config, error) {
	return Config{
		API: APIConfig{
			Port:           os.Getenv("DISTRIBUTOR_PORT"),
			PSK:            os.Getenv("DISTRIBUTOR_PSK"),
			PSKAllowOrigin: os.Getenv("DISTRIBUTOR_PSK_ALLOW_ORIGIN"),
			LogFile:        os.Getenv("DISTRIBUTOR_LOG_FILE"),
		},
	}, nil
}
