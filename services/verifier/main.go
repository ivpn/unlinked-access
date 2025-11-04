package main

import (
	"log"
	"os"

	"ivpn.net/auth/services/verifier/client"
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

	verifier, err := client.NewVerifierAWS(cfg)
	if err != nil {
		log.Println(err)
	}

	service, err := service.New(cfg, db, verifier)
	if err != nil {
		log.Println(err)
	}

	// Handle subcommands (default = serve)
	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "sync":
			if err := service.SyncManifest(); err != nil {
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
