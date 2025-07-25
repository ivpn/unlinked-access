package config

import (
	"os"
	"strings"
)

type RedisConfig struct {
	Addr             string
	Addrs            []string
	MasterName       string
	Username         string
	Password         string
	FailoverUsername string
	FailoverPassword string
}

type Config struct {
	Redis RedisConfig
}

func New() (Config, error) {
	redisAddrs := strings.Split(os.Getenv("REDIS_ADDRESSES"), ",")

	return Config{
		Redis: RedisConfig{
			Addr:             os.Getenv("REDIS_ADDR"),
			Addrs:            redisAddrs,
			MasterName:       os.Getenv("REDIS_MASTER_NAME"),
			Username:         os.Getenv("REDIS_USERNAME"),
			Password:         os.Getenv("REDIS_PASSWORD"),
			FailoverUsername: os.Getenv("REDIS_FAILOVER_USERNAME"),
			FailoverPassword: os.Getenv("REDIS_FAILOVER_PASSWORD"),
		},
	}, nil
}
