package service

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
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
const STALE_DAYS = 1

type Store interface {
	GetAccounts() ([]*model.Account, error)
	GetAccountsMock(int) ([]*model.Account, error)
	CreateAccountsMock(int) error
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

	err = gocron.Every(1).Minute().Do(RemoveStaleManifests(BASE_PATH))
	if err != nil {
		log.Printf("error scheduling manifest cleanup: %v", err)
	}

	// Start all the pending jobs
	<-gocron.Start()

	return err
}

func (s *Service) Generate() error {
	log.Println("generating manifest...")

	m, err := s.CreateManifest()
	if err != nil {
		log.Printf("error generating manifest: %v", err)
		return err
	}

	err = SignManifest(m)
	if err != nil {
		log.Printf("error signing manifest: %v", err)
		return err
	}

	err = SaveManifest(m)
	if err != nil {
		log.Printf("error saving manifest: %v", err)
		return err
	}

	log.Printf("generated manifest: %+v", m.ID)

	return nil
}

func (s *Service) CreateManifest() (*model.Manifest, error) {
	log.Println("creating manifest...")

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
		accountIDHash := sha512.Sum512([]byte(account.ID))
		token, err := s.Token.GenerateToken(string(accountIDHash[:]))
		if err != nil {
			log.Printf("error generating token for account %s: %v", account.ID, err)
			continue
		}

		tokenHash := sha256.Sum256([]byte(token))

		subscriptions[i] = model.Subscription{
			TokenHash:   string(tokenHash[:]),
			IsActive:    account.IsActive,
			ActiveUntil: account.ActiveUntil,
			Tier:        account.Product,
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

func SignManifest(m *model.Manifest) error {
	// TODO: Implement HSM signing
	log.Println("signing manifest...")

	data, err := json.Marshal(m)
	if err != nil {
		log.Println("error marshaling manifest for signing:", err)
		return err
	}

	hash := sha256.Sum256(data)
	m.Signature = base64.StdEncoding.EncodeToString(hash[:])

	log.Printf("manifest signed: %s", m.Signature)

	return nil
}

func RemoveStaleManifests(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read manifest directory: %w", err)
	}

	// Calculate the cutoff time for stale manifests
	cutoff := time.Now().AddDate(0, 0, -STALE_DAYS)

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
