package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jasonlvhit/gocron"
	"golang.org/x/time/rate"
	"ivpn.net/auth/services/generator/config"
	"ivpn.net/auth/services/generator/model"
)

// Configuration
const (
	currentManifest = "current.json"
	basePath        = "/app/data"
	staleDays       = 1
	burst           = 1 // allowable burst above targetTPS
	workerCount     = 4 // number of signing goroutines
)

type Store interface {
	GetAccounts() ([]*model.Account, error)
	CreateAccountsMock(int) error
}

type TokenClient interface {
	GenerateToken(string) (string, error)
}

type Service struct {
	Cfg   config.Config
	Store Store
	Token TokenClient
}

func New(cfg config.Config, store Store, tokenClient TokenClient) *Service {
	return &Service{
		Cfg:   cfg,
		Store: store,
		Token: tokenClient,
	}
}

func (s *Service) Start() error {
	log.Println("generator service started")

	if s.Cfg.Service.SampleData {
		err := s.Store.CreateAccountsMock(25)
		if err != nil {
			log.Printf("error creating mock accounts: %v", err)
			return err
		}
	}

	err := gocron.Every(1).Hour().Do(s.Generate)
	if err != nil {
		log.Printf("error scheduling manifest generation: %v", err)
	}

	err = gocron.Every(1).Hour().Do(s.RemoveStaleManifests)
	if err != nil {
		log.Printf("error scheduling stale manifest removal: %v", err)
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

	err = s.SignManifest(m)
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
	accounts, err := s.Store.GetAccounts()
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

	ctx := context.Background()
	limiter := rate.NewLimiter(rate.Limit(s.Cfg.Service.TPS), burst)
	jobs := make(chan *model.Account, len(accounts))
	results := make(chan model.Subscription, len(accounts))
	var wg sync.WaitGroup

	// Start workers
	for w := range workerCount {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for account := range jobs {
				// Wait for limiter before each Sign call
				if err := limiter.Wait(ctx); err != nil {
					log.Printf("[worker %d] limiter error: %v", workerID, err)
					continue
				}

				// Generate token for account ID
				token, err := s.Token.GenerateToken(account.ID)
				if err != nil {
					continue
				}

				// Hash the token
				tokenHash := sha256.Sum256([]byte(token))

				// Set ActiveUntil to 12:00 (UTC) of the day from account.ActiveUntil
				accountUntil := account.ActiveUntil.UTC()
				roundedUntil := time.Date(accountUntil.Year(), accountUntil.Month(), accountUntil.Day(), 12, 0, 0, 0, time.UTC)

				var tier string
				switch account.Product {
				case "IVPN Standard":
					tier = "IVPN Tier 1"
				case "IVPN Pro":
					tier = "IVPN Tier 3"
				default:
					tier = account.Product
				}

				results <- model.Subscription{
					TokenHash:   base64.StdEncoding.EncodeToString(tokenHash[:]),
					IsActive:    account.IsActive,
					ActiveUntil: roundedUntil,
					Tier:        tier,
				}
			}
		}(w)
	}

	start := time.Now()

	// Send jobs
	for _, account := range accounts {
		jobs <- account
	}
	close(jobs)

	// Wait for workers to finish
	wg.Wait()
	close(results)

	// Collect results
	signedSubs := make([]model.Subscription, 0, len(accounts))
	for sub := range results {
		signedSubs = append(signedSubs, sub)
	}

	// Randomize order
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(signedSubs), func(i, j int) {
		signedSubs[i], signedSubs[j] = signedSubs[j], signedSubs[i]
	})

	log.Printf("signed %d subscriptions in %s with %d workers (limit: %d TPS)\n", len(signedSubs), time.Since(start), workerCount, s.Cfg.Service.TPS)

	return signedSubs, nil
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
	basePath := basePath
	timestampFile := fmt.Sprintf("%s/%s.json", basePath, timestamp)
	currentFile := fmt.Sprintf("%s/%s", basePath, currentManifest)

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

func (s *Service) SignManifest(m *model.Manifest) error {
	log.Println("signing manifest...")

	data, err := json.Marshal(m)
	if err != nil {
		log.Println("error marshaling manifest for signing:", err)
		return err
	}

	digest := sha256.Sum256(data)
	digestBase64 := base64.StdEncoding.EncodeToString(digest[:])

	// Generate signature for manifest hash
	signature, err := s.Token.GenerateToken(digestBase64)
	if err != nil {
		log.Println("error generating token for manifest hash:", err)
		return err
	}

	m.Signature = signature

	log.Printf("manifest signed: %s", m.Signature)

	return nil
}

func (s *Service) RemoveStaleManifests() error {
	log.Println("deleting stale manifests...")

	files, err := os.ReadDir(basePath)
	if err != nil {
		return fmt.Errorf("failed to read manifest directory: %w", err)
	}

	// Calculate the cutoff time for stale manifests
	cutoff := time.Now().AddDate(0, 0, -staleDays)

	for _, file := range files {
		name := file.Name()

		// Skip non-JSON files and current manifest
		if file.IsDir() || !strings.HasSuffix(name, ".json") || name == currentManifest {
			continue
		}

		// Try parsing timestamp from filename (e.g., 2025-07-15T10-04-00.json)
		timestampStr := strings.TrimSuffix(name, ".json")
		timestamp, err := time.Parse("2006-01-02T15-04-05", timestampStr)
		if err != nil {
			continue
		}

		if timestamp.Before(cutoff) {
			fullPath := filepath.Join(basePath, name)
			if err := os.Remove(fullPath); err != nil {
				fmt.Printf("failed to delete old manifest %s: %v\n", name, err)
			} else {
				fmt.Printf("deleted old manifest: %s\n", name)
			}
		}
	}

	return nil
}
