package service

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jasonlvhit/gocron"
	"ivpn.net/auth/services/generator/model"
)

type Store interface {
	GetAccounts() ([]*model.Account, error)
	GetAccountsMock(int) ([]*model.Account, error)
}

type TokenClient interface {
	GenerateToken(string) (string, error)
}

type Service struct {
	Store Store
	Token TokenClient
}

func New(store Store, tokenClient TokenClient) *Service {
	return &Service{
		Store: store,
		Token: tokenClient,
	}
}

func (s *Service) Start() error {
	log.Println("generator service started")

	err := gocron.Every(1).Minute().Do(s.Generate)
	if err != nil {
		log.Printf("error scheduling account retrieval: %v", err)
	}

	// Start all the pending jobs
	<-gocron.Start()

	return err
}

func (s *Service) Generate() error {
	log.Println("generating manifest...")

	manifest, err := s.GenerateManifest()
	if err != nil {
		log.Printf("error generating manifest: %v", err)
		return err
	}

	log.Printf("generated manifest: %+v", manifest.ID)

	return nil
}

func (s *Service) GenerateManifest() (*model.Manifest, error) {
	log.Println("generating metadata...")

	subs, err := s.GenerateSubscriptions()
	if err != nil {
		log.Printf("error generating subscriptions: %v", err)
		return nil, err
	}

	manifest := &model.Manifest{
		ID:            uuid.New().String(),
		CreatedAt:     time.Now(),
		ValidUntil:    time.Now().Add(3 * time.Hour),
		Subscriptions: subs,
		Signature:     "",
	}

	return manifest, nil
}

func (s *Service) GetAccounts() ([]*model.Account, error) {
	accounts, err := s.Store.GetAccountsMock(10)
	if err != nil {
		log.Printf("error fetching accounts: %v", err)
		return nil, err
	}

	return accounts, nil
}

func (s *Service) GenerateSubscriptions() ([]model.Subscription, error) {
	log.Println("generating subscriptions...")

	accounts, err := s.GetAccounts()
	if err != nil {
		log.Printf("error fetching accounts: %v", err)
		return nil, err
	}

	subscriptions := make([]model.Subscription, len(accounts))
	for i, account := range accounts {
		token, err := s.Token.GenerateToken(account.ID)
		if err != nil {
			log.Printf("error generating token for account %s: %v", account.ID, err)
			continue
		}

		sha256Token := sha256.Sum256([]byte(token))
		token = string(sha256Token[:])

		subscriptions[i] = model.Subscription{
			TokenHash:   token,
			IsActive:    account.IsActive,
			ActiveUntil: account.ActiveUntil,
		}
	}

	return subscriptions, nil
}

func (s *Service) SaveManifest(manifest *model.Manifest) error {
	log.Printf("saving manifest: %s", manifest.ID)

	// Marshal the manifest to JSON
	jsonData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		log.Println("error marshalling to JSON:", err)
		return err
	}

	timestamp := time.Now().Format("2006-01-02T15-04-05")
	basePath := "/app/data"
	timestampFilename := fmt.Sprintf("%s/%s.json", basePath, timestamp)
	currentFilename := fmt.Sprintf("%s/current-manifest.json", basePath)

	// Write both files
	if err := os.WriteFile(timestampFilename, jsonData, 0644); err != nil {
		log.Println("error writing timestamped file:", err)
		return err
	}
	if err := os.WriteFile(currentFilename, jsonData, 0644); err != nil {
		log.Println("error writing current-manifest file:", err)
		return err
	}

	log.Println("manifest saved:")
	log.Println(" -", timestampFilename)
	log.Println(" -", currentFilename)

	return nil
}
