package service

import (
	"context"
	"time"

	"ivpn.net/auth/services/preauth/config"
	"ivpn.net/auth/services/preauth/model"
)

type Cache interface {
	Set(context.Context, string, any, time.Duration) error
	Get(context.Context, string) (string, error)
	Del(context.Context, string) error
	Incr(context.Context, string, time.Duration) error
}

type Service struct {
	Cfg   config.Config
	Cache Cache
}

func New(cfg config.Config, cache Cache) *Service {
	return &Service{
		Cfg:   cfg,
		Cache: cache,
	}
}

func (s *Service) GetPreAuth(ID string) (model.PreAuth, error) {
	return model.PreAuth{}, nil
}

func (s *Service) AddPreAuth(accountId string) error {
	return nil
}
