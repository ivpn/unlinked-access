package main

import (
	"log"

	"ivpn.net/auth/services/distributor/api"
	"ivpn.net/auth/services/distributor/config"
	"ivpn.net/auth/services/distributor/service"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Println(err)
	}

	service := service.New(cfg)

	err = api.Start(cfg.API, service)
	if err != nil {
		log.Println(err)
	}
}
