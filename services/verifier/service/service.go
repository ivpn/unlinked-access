package service

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

	err = VerifyManifest(manifest)
	if err != nil {
		log.Printf("manifest verification failed: %v", err)
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

func VerifyManifest(m model.Manifest) error {
	// TODO: Implement HSM verification
	log.Printf("verifying manifest: %v", m.ID)

	signature := m.Signature
	m.Signature = ""

	data, err := json.Marshal(m)
	if err != nil {
		log.Println("error marshaling manifest for signing:", err)
		return err
	}

	hash := sha256.Sum256(data)
	hashString := base64.StdEncoding.EncodeToString(hash[:])

	if hashString != signature {
		log.Printf("manifest signature does not match: %v != %v", hashString, signature)
		return fmt.Errorf("invalid manifest signature")
	}

	log.Println("manifest signature OK")

	return nil
}
