package model

import "time"

type Manifest struct {
	ID            string         `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	ValidUntil    time.Time      `json:"valid_until"`
	Subscriptions []Subscription `json:"subscriptions"`
	Signature     string         `json:"signature"`
}
