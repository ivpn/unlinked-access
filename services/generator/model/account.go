package model

import "time"

type Account struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	ActiveUntil time.Time `json:"active_until"`
	Product     string    `json:"product"`
}
