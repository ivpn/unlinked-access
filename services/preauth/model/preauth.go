package model

import "time"

type PreAuth struct {
	ID          string    `json:"id"`
	TokenHash   string    `json:"token_hash"`
	ActiveUntil time.Time `json:"active_until"`
	Tier        string    `json:"tier"`
}
