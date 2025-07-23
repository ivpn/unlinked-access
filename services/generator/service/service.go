package service

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jasonlvhit/gocron"
	"ivpn.net/auth/services/generator/model"
)

const CURRENT_MANIFEST = "current.json"
const BASE_PATH = "/app/data"

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
		log.Printf("error scheduling manifest generation: %v", err)
	}

	err = gocron.Every(1).Minute().Do(CleanupOldManifests(BASE_PATH))
	if err != nil {
		log.Printf("error scheduling manifest cleanup: %v", err)
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

func SaveManifest(m *model.Manifest) error {
	log.Printf("saving manifest: %s", m.ID)

	// Marshal the manifest to JSON
	jsonData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Println("error converting manifest to JSON:", err)
		return err
	}

	timestamp := time.Now().Format("2006-01-02T15-04-05")
	basePath := BASE_PATH
	timestampFile := fmt.Sprintf("%s/%s.json", basePath, timestamp)
	currentFile := fmt.Sprintf("%s/%s", basePath, CURRENT_MANIFEST)

	// Write both files
	if err := os.WriteFile(timestampFile, jsonData, 0600); err != nil {
		log.Println("error writing timestamp file:", err)
		return err
	}
	if err := os.WriteFile(currentFile, jsonData, 0600); err != nil {
		log.Println("error writing current file:", err)
		return err
	}

	log.Println("manifest saved:")
	log.Println(" -", timestampFile)
	log.Println(" -", currentFile)

	return nil
}

// Delete all JSON manifest files older than 7 days
func CleanupOldManifests(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read manifest directory: %w", err)
	}

	// 7-day cutoff
	cutoff := time.Now().AddDate(0, 0, -7)

	for _, file := range files {
		name := file.Name()

		// Skip non-JSON files and current manifest
		if file.IsDir() || !strings.HasSuffix(name, ".json") || name == CURRENT_MANIFEST {
			continue
		}

		// Try parsing timestamp from filename (e.g., 2025-07-15T10-04-00.json)
		timestampStr := strings.TrimSuffix(name, ".json")
		timestamp, err := time.Parse("2006-01-02T15-04-05", timestampStr)
		if err != nil {
			continue
		}

		if timestamp.Before(cutoff) {
			fullPath := filepath.Join(dir, name)
			if err := os.Remove(fullPath); err != nil {
				fmt.Printf("failed to delete old manifest %s: %v\n", name, err)
			} else {
				fmt.Printf("deleted old manifest: %s\n", name)
			}
		}
	}

	return nil
}
