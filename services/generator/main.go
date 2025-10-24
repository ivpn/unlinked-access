package main

import (
	"log"
	"os"

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

	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Println(err)
	}

	tokenClient, err := client.New(cfg.TokenServer)
	if err != nil {
		log.Println(err)
	}

	service := service.New(cfg, db, tokenClient)

	// Handle subcommands (default = serve)
	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "sync":
			if err := service.Generate(); err != nil {
				log.Println(err)
			}
			return

		case "serve":
			// continue to service.Start below
			break

		default:
			return
		}
	}

	err = service.Start()
	if err != nil {
		log.Println(err)
	}
}
