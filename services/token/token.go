package main

import "ivpn.net/auth/services/token/model"

type HSMClient interface {
	Token(input string, ttlMinutes int) (*model.HSMToken, error)
}

type TokenService struct {
	HSMClient HSMClient
}

func New(hsmClient HSMClient) *TokenService {
	return &TokenService{
		HSMClient: hsmClient,
	}
}

func (s *TokenService) GenerateToken(input string, ttlMinutes int) (*model.HSMToken, error) {
	return s.HSMClient.Token(input, ttlMinutes)
}
