package main

import (
	"log"

	"ivpn.net/auth/services/verifier/config"
	"ivpn.net/auth/services/verifier/repository"
	"ivpn.net/auth/services/verifier/service"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Println(err)
	}

	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Println(err)
	}

	service, err := service.New(cfg, db)
	if err != nil {
		log.Println(err)
	}

	err = service.Start()
	if err != nil {
		log.Println(err)
	}
}
