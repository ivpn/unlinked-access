package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jasonlvhit/gocron"
	"ivpn.net/auth/services/verifier/client/http"
	"ivpn.net/auth/services/verifier/config"
	"ivpn.net/auth/services/verifier/model"
)

type Store interface {
	GetSubscriptions() ([]model.Subscription, error)
	UpdateSubscriptions([]model.Subscription) error
	GetLatestManifestLog() (model.ManifestLog, error)
	AddManifestLog(model.ManifestLog) error
	CleanupManifestLogs() error
}

type Verifier interface {
	Verify(signature string, data []byte) error
	Authenticate() error
	IsAuthError(err error) bool
}

type Service struct {
	Cfg      config.Config
	Stores   []Store
	Http     http.Http
	Verifier Verifier
}

func New(cfg config.Config, stores []Store, verifier Verifier) (*Service, error) {
	return &Service{
		Cfg:      cfg,
		Stores:   stores,
		Verifier: verifier,
		Http: http.Http{
			Cfg: cfg.API,
		},
	}, nil
}

func (s *Service) Start() error {
	log.Println("verifier service started")

	err := gocron.Every(1).Hour().Do(s.SyncManifest)
	if err != nil {
		log.Printf("error syncing manifest: %v", err)
	}

	err = gocron.Every(1).Day().Do(s.CleanupManifestLogs)
	if err != nil {
		log.Printf("error cleaning up manifest logs: %v", err)
	}

	// Start all the pending jobs
	<-gocron.Start()

	return err
}

func (s *Service) SyncManifest() error {
	m, err := s.GetManifest()
	if err != nil {
		return err
	}

	manifestLog := s.CreateManifestLog(m)

	err = s.VerifyManifest(m)
	if err != nil {
		err = s.SaveManifestLog(m, manifestLog)
		if err != nil {
			return err
		}

		return err
	}

	manifestLog.SignatureValid = true

	err = s.UpdateSubscriptions(m)
	if err != nil {
		err = s.SaveManifestLog(m, manifestLog)
		if err != nil {
			return err
		}

		return err
	}

	manifestLog.Status = "success"

	err = s.SaveManifestLog(m, manifestLog)
	if err != nil {
		return err
	}

	log.Printf("manifest synced successfully: %v", m.ID)

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

func (s *Service) VerifyManifest(m model.Manifest) error {
	log.Printf("verifying manifest: %v", m.ID)

	if m.ValidUntil.Before(time.Now()) {
		log.Printf("manifest is expired: %v", m.ValidUntil)
		return fmt.Errorf("manifest is expired")
	}

	signature := m.Signature
	m.Signature = ""

	data, err := json.Marshal(m)
	if err != nil {
		log.Println("error marshaling manifest for signing:", err)
		return err
	}

	err = s.Verifier.Verify(signature, data)
	if err != nil {
		if s.Verifier.IsAuthError(err) {
			log.Println("re-authenticating verifier session...")
			if authErr := s.Verifier.Authenticate(); authErr == nil {
				err = s.Verifier.Verify(signature, data)
				if err != nil {
					log.Println(err)
					return err
				}
				return nil
			}
		}

		log.Println("error verifying manifest signature:", err)
		return err
	}

	return nil
}

func (s *Service) VerifyManifestVersion(m model.Manifest, store Store) error {
	log.Printf("verifying manifest version: %v", m.Version)

	latestLog, err := store.GetLatestManifestLog()
	if err != nil {
		log.Printf("error fetching latest manifest log: %v", err)
		return err
	}

	if m.Version <= latestLog.Version {
		return fmt.Errorf("manifest version is not newer than the latest log")
	}

	return nil
}

func (s *Service) UpdateSubscriptions(m model.Manifest) error {
	var lastErr error
	for _, store := range s.Stores {
		if err := s.VerifyManifestVersion(m, store); err != nil {
			log.Printf("error verifying manifest version: %v", err)
			lastErr = err
			continue
		}

		subs, err := store.GetSubscriptions()
		if err != nil {
			log.Printf("error fetching subscriptions from store: %v", err)
			lastErr = err
			continue
		}

		var updatedSubs []model.Subscription
		for _, sub := range subs {
			updatedSub, err := UpdateSubscriptionFromManifest(sub, m.Subscriptions)
			if err != nil {
				continue
			}
			updatedSubs = append(updatedSubs, updatedSub)
		}

		if err := store.UpdateSubscriptions(updatedSubs); err != nil {
			log.Printf("error saving updated subscriptions to store: %v", err)
			lastErr = err
			continue
		}
		log.Printf("updated %d subscriptions in store", len(updatedSubs))
	}

	return lastErr
}

func (s *Service) CreateManifestLog(m model.Manifest) model.ManifestLog {
	return model.ManifestLog{
		ID:             uuid.New().String(),
		Version:        m.Version,
		CreatedAt:      time.Now(),
		SignatureValid: false,
		Status:         "failed",
	}
}

func (s *Service) SaveManifestLog(m model.Manifest, manifestLog model.ManifestLog) error {
	for _, store := range s.Stores {
		if err := store.AddManifestLog(manifestLog); err != nil {
			log.Printf("error adding manifest log entry to store: %v", err)
			return err
		}
	}
	return nil
}

func (s *Service) CleanupManifestLogs() error {
	for _, store := range s.Stores {
		if err := store.CleanupManifestLogs(); err != nil {
			log.Printf("error removing expired manifest logs from store: %v", err)
			return err
		}
	}
	return nil
}

func UpdateSubscriptionFromManifest(sub model.Subscription, manifestSubs []model.Subscription) (model.Subscription, error) {
	for _, s := range manifestSubs {
		if sub.TokenHash == s.TokenHash {
			sub.ActiveUntil = s.ActiveUntil
			sub.Tier = s.Tier
			return sub, nil
		}
	}

	return model.Subscription{}, fmt.Errorf("subscription %s not found in manifest", sub.ID)
}
