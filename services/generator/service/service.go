package service

import (
	"log"

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
	// Initialize the service
	log.Println("generator service started")
	return nil
}
