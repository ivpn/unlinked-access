package model

import "time"

type Account struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	IsActive    bool      `json:"is_active"`
	ActiveUntil time.Time `json:"active_until"`
	Product     string    `json:"product"`
}
