package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jasonlvhit/gocron"
	"ivpn.net/auth/services/verifier/client/http"
	"ivpn.net/auth/services/verifier/config"
	"ivpn.net/auth/services/verifier/model"
)

type Store interface {
	GetSubscriptions() ([]model.Subscription, error)
	UpdateSubscriptions([]model.Subscription) error
}

type Verifier interface {
	Verify(signature string, data []byte) error
}

type Service struct {
	Cfg      config.Config
	Store    Store
	Http     http.Http
	Verifier Verifier
}

func New(cfg config.Config, store Store, verifier Verifier) (*Service, error) {
	return &Service{
		Cfg:      cfg,
		Store:    store,
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

	// Start all the pending jobs
	<-gocron.Start()

	return err
}

func (s *Service) SyncManifest() error {
	log.Println("syncing manifest...")
	m, err := s.GetManifest()
	if err != nil {
		return err
	}

	err = s.VerifyManifest(m)
	if err != nil {
		return err
	}

	err = s.UpdateSubscriptions(m)
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
		log.Println("error verifying manifest signature:", err)
		return err
	}

	return nil
}

func (s *Service) UpdateSubscriptions(m model.Manifest) error {
	subs, err := s.Store.GetSubscriptions()
	if err != nil {
		log.Printf("error fetching subscriptions: %v", err)
		return err
	}

	for i, sub := range subs {
		updatedSub, err := UpdateSubscriptionFromManifest(sub, m.Subscriptions)
		if err != nil {
			log.Printf("error updating subscription: %v", err)
			continue
		}

		subs[i] = updatedSub
	}

	err = s.Store.UpdateSubscriptions(subs)
	if err != nil {
		log.Printf("error saving updated subscriptions: %v", err)
		return err
	}

	return nil
}

func UpdateSubscriptionFromManifest(sub model.Subscription, manifestSubs []model.Subscription) (model.Subscription, error) {
	for _, s := range manifestSubs {
		if sub.TokenHash == s.TokenHash {
			sub.IsActive = s.IsActive
			sub.ActiveUntil = s.ActiveUntil
			sub.Tier = s.Tier
			return sub, nil
		}
	}

	return model.Subscription{}, fmt.Errorf("subscription with TokenHash %s not found", sub.TokenHash)
}
