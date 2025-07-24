package main

import (
	"log"

	"ivpn.net/auth/services/distributor/api"
	"ivpn.net/auth/services/distributor/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Println(err)
	}

	err = api.Start(cfg.API)
	if err != nil {
		log.Println(err)
	}
}
