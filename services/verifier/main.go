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

	db, err := repository.NewDB(cfg.DB)
	if err != nil {
		log.Println(err)
	}

	service := service.New(cfg, db)
	err = service.Start()
	if err != nil {
		log.Println(err)
	}
}
