package model

import "time"

type Subscription struct {
	TokenHash   string    `json:"h"`
	ActiveUntil time.Time `json:"u"`
	Tier        string    `json:"t"`
}
