package service

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"ivpn.net/auth/services/preauth/config"
	"ivpn.net/auth/services/preauth/model"
)

type Cache interface {
	Set(context.Context, string, any, time.Duration) error
	Get(context.Context, string) (string, error)
	Del(context.Context, string) error
	Incr(context.Context, string, time.Duration) error
}

type TokenClient interface {
	GenerateToken(string) (string, error)
}

type Service struct {
	Cfg   config.Config
	Cache Cache
	Token TokenClient
}

func New(cfg config.Config, cache Cache, token TokenClient) *Service {
	return &Service{
		Cfg:   cfg,
		Cache: cache,
		Token: token,
	}
}

func (s *Service) GetPreAuth(ctx context.Context, ID string) (model.PreAuth, error) {
	// Retrieve data from Cache
	val, err := s.Cache.Get(ctx, "preauth_"+ID)
	if err != nil {
		log.Println("failed to get pre-auth from cache:", err)
		return model.PreAuth{}, err
	}

	// Unmarshal the JSON into a struct
	var retrieved model.PreAuth
	err = json.Unmarshal([]byte(val), &retrieved)
	if err != nil {
		log.Println("failed to unmarshal pre-auth from cache:", err)
		return model.PreAuth{}, err
	}

	return retrieved, nil
}

func (s *Service) AddPreAuth(ctx context.Context, accountId string) error {
	// Generate token
	accountIDHash := sha512.Sum512([]byte(accountId))
	token, err := s.Token.GenerateToken(string(accountIDHash[:]))
	if err != nil {
		log.Println("failed to generate token:", err)
		return err
	}

	// Create an instance of PreAuth
	tokenHash := sha256.Sum256([]byte(token))
	pa := model.PreAuth{
		ID:        uuid.New().String(),
		TokenHash: string(tokenHash[:]),
	}

	// Marshal the struct to JSON
	data, err := json.Marshal(pa)
	if err != nil {
		log.Println("failed to marshal pre-auth to JSON:", err)
		return err
	}

	// Set in Redis
	err = s.Cache.Set(ctx, "preauth_"+pa.ID, string(data), s.Cfg.API.PreauthTTL)
	if err != nil {
		log.Println("failed to set pre-auth in cache:", err)
		return err
	}

	return nil
}
