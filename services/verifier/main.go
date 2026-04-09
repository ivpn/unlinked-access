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
		log.Fatal(err)
	}

	var stores []service.Store

	if cfg.DB.Host != "" {
		db, err := repository.NewDB(cfg)
		if err != nil {
			log.Fatal(err)
		}
		stores = append(stores, db)
	}

	if cfg.NoSQLDB.Host != "" {
		mdb, err := repository.NewMongoDB(cfg)
		if err != nil {
			log.Fatal(err)
		}
		stores = append(stores, mdb)
	}

	if len(stores) == 0 {
		log.Fatal("no stores configured: set CLIENT_DB_HOST and/or CLIENT_DB_NOSQL_HOST")
	}

	verifier, err := client.NewVerifierFortanix(cfg)
	if err != nil {
		log.Fatal(err)
	}

	svc, err := service.New(cfg, stores, verifier)
	if err != nil {
		log.Fatal(err)
	}

	// Handle subcommands (default = serve)
	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "sync":
			if err := svc.SyncManifest(); err != nil {
				log.Println(err)
			}
			return

		case "serve":
			// continue to svc.Start below
			break

		default:
			return
		}
	}

	if err := svc.Start(); err != nil {
		log.Println(err)
	}
}
