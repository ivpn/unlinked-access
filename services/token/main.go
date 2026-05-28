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
		log.Fatal(err)
	}

	signer, err := client.NewSignerFortanix(cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := service.New(signer, cfg)
	if err = server.Start(); err != nil {
		log.Fatal(err)
	}
}
