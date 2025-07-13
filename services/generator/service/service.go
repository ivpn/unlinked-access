package service

import (
	"log"

	"github.com/jasonlvhit/gocron"
	"ivpn.net/auth/services/generator/model"
)

type Store interface {
	GetAccounts() ([]*model.Account, error)
	GetAccountsMock(count int) ([]*model.Account, error)
}

type Service struct {
	Store Store
}

func New(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) Start() error {
	log.Println("generator service started")

	err := gocron.Every(1).Minute().Do(s.Generate)
	if err != nil {
		log.Printf("error scheduling account retrieval: %v", err)
	}

	return err
}

func (s *Service) Generate() error {
	log.Println("generating manifest...")

	accounts, err := s.Store.GetAccounts()
	if err != nil {
		log.Printf("error fetching accounts: %v", err)
		return err
	}

	// Process the accounts
	for _, account := range accounts {
		log.Printf("processing account: %v", account.ID)
	}

	return nil
}

func (s *Service) GetAccounts() ([]*model.Account, error) {
	accounts, err := s.Store.GetAccountsMock(10)
	if err != nil {
		log.Printf("error fetching accounts: %v", err)
		return nil, err
	}

	return accounts, nil
}
