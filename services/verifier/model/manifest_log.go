package model

import "time"

type ManifestLog struct {
	ID             string    `json:"id"`
	Version        int       `json:"version"`
	CreatedAt      time.Time `json:"created_at"`
	SignatureValid bool      `json:"signature_valid"`
	Status         string    `json:"status"`
}
