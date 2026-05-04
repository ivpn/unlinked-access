package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID `json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty"`
	TokenHash   string    `json:"h" gorm:"column:token_hash" bson:"token_hash"`
	IsActive    bool      `json:"a" gorm:"column:is_active" bson:"is_active"`
	ActiveUntil time.Time `json:"u" gorm:"column:active_until" bson:"active_until"`
	Tier        string    `json:"t" gorm:"column:tier" bson:"tier"`
}
