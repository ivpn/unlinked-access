package model

import "time"

type Subscription struct {
	ID          string    `json:"id"`
	TokenHash   string    `json:"h"`
	IsActive    bool      `json:"a"`
	ActiveUntil time.Time `json:"u"`
	Tier        string    `json:"t"`
}
