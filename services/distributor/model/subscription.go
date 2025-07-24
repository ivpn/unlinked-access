package model

import "time"

type Subscription struct {
	TokenHash   string    `json:"token_hash"`
	IsActive    bool      `json:"is_active"`
	ActiveUntil time.Time `json:"active_until"`
}
