package main

import (
	"log"

	"ivpn.net/auth/services/generator/client"
	"ivpn.net/auth/services/generator/config"
	"ivpn.net/auth/services/generator/repository"
	"ivpn.net/auth/services/generator/service"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Println(err)
	}

	db, err := repository.NewDB(cfg.DB)
	if err != nil {
		log.Println(err)
	}

	tokenClient, err := client.New(cfg.TokenServer)
	if err != nil {
		log.Println(err)
	}

	service := service.New(db, tokenClient)
	err = service.Start()
	if err != nil {
		log.Println(err)
	}
}
