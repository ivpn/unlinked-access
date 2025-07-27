package service

import (
	"log"

	"github.com/jasonlvhit/gocron"
	"ivpn.net/auth/services/verifier/client/http"
	"ivpn.net/auth/services/verifier/config"
	"ivpn.net/auth/services/verifier/model"
)

type Service struct {
	Http http.Http
}

func New(cfg config.Config) *Service {
	return &Service{
		Http: http.Http{
			Cfg: cfg.API,
		},
	}
}

func (s *Service) Start() error {
	log.Println("verifier service started")

	err := gocron.Every(1).Minute().Do(s.SyncManifest)
	if err != nil {
		log.Printf("error fetching manifest: %v", err)
	}

	// Start all the pending jobs
	<-gocron.Start()

	return err
}

func (s *Service) SyncManifest() error {
	log.Println("syncing manifest...")
	manifest, err := s.GetManifest()
	if err != nil {
		log.Printf("error syncing manifest: %v", err)
		return err
	}

	log.Printf("manifest synced successfully: %v", manifest.ID)

	return nil
}

func (s *Service) GetManifest() (model.Manifest, error) {
	manifest, err := s.Http.GetManifest()
	if err != nil {
		log.Printf("error fetching manifest: %v", err)
		return model.Manifest{}, err
	}

	return manifest, nil
}
