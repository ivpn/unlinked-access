package model

type Session struct {
	ID        string `json:"id"`
	Token     string `json:"token"`
	PreAuthID string `json:"preauth_id"`
}
