package model

import "time"

type Subscription struct {
	ID          string    `json:"id,omitempty" gorm:"column:id"`
	TokenHash   string    `json:"h" gorm:"column:token_hash"`
	IsActive    bool      `json:"a" gorm:"column:is_active"`
	ActiveUntil time.Time `json:"u" gorm:"column:active_until"`
	Tier        string    `json:"t" gorm:"column:tier"`
}
