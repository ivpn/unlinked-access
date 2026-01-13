package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"ivpn.net/auth/services/preauth/client/http"
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
	Http  *http.Http
}

func New(cfg config.Config, cache Cache, token TokenClient) *Service {
	return &Service{
		Cfg:   cfg,
		Cache: cache,
		Token: token,
		Http:  http.New(cfg.API),
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

func (s *Service) AddPreAuth(ctx context.Context, accountId string, isActive bool, activeUntil time.Time, tier string) (string, error) {
	// Generate token
	token, err := s.Token.GenerateToken(accountId)
	if err != nil {
		log.Println("failed to generate token:", err)
		return "", err
	}

	// Create an instance of PreAuth
	tokenHash := sha256.Sum256([]byte(token))
	pa := model.PreAuth{
		ID:          uuid.New().String(),
		TokenHash:   base64.StdEncoding.EncodeToString(tokenHash[:]),
		IsActive:    isActive,
		ActiveUntil: activeUntil,
		Tier:        tier,
	}

	// Marshal the struct to JSON
	data, err := json.Marshal(pa)
	if err != nil {
		log.Println("failed to marshal pre-auth to JSON:", err)
		return "", err
	}

	// Set in Redis
	err = s.Cache.Set(ctx, "preauth_"+pa.ID, string(data), s.Cfg.API.PreauthTTL)
	if err != nil {
		log.Println("failed to set pre-auth in cache:", err)
		return "", err
	}

	// Create an instance of Session
	session := model.Session{
		ID:        uuid.New().String(),
		Token:     token,
		PreAuthID: pa.ID,
	}

	// Post session to webhooks
	for i, url := range s.Cfg.API.SessionURLs {
		psk := s.Cfg.API.SessionPSKs[i]
		err = s.Http.PostSession(session, url, psk)
		if err != nil {
			log.Println("failed to post session to webhook:", err)
			return "", err
		}
	}

	return session.ID, nil
}
