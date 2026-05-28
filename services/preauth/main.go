package main

import (
	"log"

	"ivpn.net/auth/services/preauth/api"
	"ivpn.net/auth/services/preauth/client"
	"ivpn.net/auth/services/preauth/config"
	"ivpn.net/auth/services/preauth/repository"
	"ivpn.net/auth/services/preauth/service"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	redis, err := repository.New(cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}

	tokenClient, err := client.New(cfg.TokenServer)
	if err != nil {
		log.Fatal(err)
	}

	service := service.New(cfg, redis, tokenClient)

	if err = api.Start(cfg.API, service); err != nil {
		log.Fatal(err)
	}
}
