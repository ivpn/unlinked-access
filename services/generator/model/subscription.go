package model

type Subscription struct {
	TokenHash   string `json:"token_hash"`
	IsActive    bool   `json:"is_active"`
	ActiveUntil string `json:"active_until"`
}
