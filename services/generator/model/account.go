package model

import "time"

type Account struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	IsActive    bool      `json:"is_active"`
	ActiveUntil time.Time `json:"active_until"`
	Product     string    `json:"product"`
	// Salt is derived from services.salt via JOIN, not stored on the accounts table
	Salt bool `json:"salt" gorm:"->;-:migration"`
}
