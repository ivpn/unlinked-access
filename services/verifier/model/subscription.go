package model

import "time"

type Subscription struct {
	TokenHash   string    `json:"h"`
	IsActive    bool      `json:"a"`
	ActiveUntil time.Time `json:"u"`
	Tier        string    `json:"t"`
}
