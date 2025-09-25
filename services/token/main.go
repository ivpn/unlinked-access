package main

import (
	"log"

	"ivpn.net/auth/services/token/client"
	"ivpn.net/auth/services/token/config"
	"ivpn.net/auth/services/token/service"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Println(err)
	}

	signer, err := client.NewAWSSigner(cfg)
	if err != nil {
		log.Println(err)
	}

	server := service.New(signer, cfg)
	err = server.Start()
	if err != nil {
		log.Println(err)
	}
}
