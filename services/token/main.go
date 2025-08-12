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

	hsm, err := client.NewHSM()
	if err != nil {
		log.Println(err)
	}

	server := service.New(hsm, cfg)
	err = server.Start()
	if err != nil {
		log.Println(err)
	}
}
