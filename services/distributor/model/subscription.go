package model

import "time"

type Subscription struct {
	TokenHash   string    `json:"t"`
	IsActive    bool      `json:"a"`
	ActiveUntil time.Time `json:"u"`
}
