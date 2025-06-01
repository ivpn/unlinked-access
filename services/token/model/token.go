package model

import "time"

type HSMToken struct {
	Token     string
	ExpiresAt time.Time
}
