package main

import (
	"log"

	"ivpn.net/auth/services/preauth/api"
	"ivpn.net/auth/services/preauth/config"
	"ivpn.net/auth/services/preauth/repository"
	"ivpn.net/auth/services/preauth/service"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Println(err)
	}

	redis, err := repository.New(cfg.Redis)
	if err != nil {
		log.Println(err)
	}

	service := service.New(cfg, redis)

	err = api.Start(cfg.API, service)
	if err != nil {
		log.Println(err)
	}
}
