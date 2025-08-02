package api

type PreauthReq struct {
	AccountID   string `json:"account_id" validate:"required"`
	IsActive    bool   `json:"is_active" validate:"required"`
	ActiveUntil string `json:"active_until" validate:"required"`
	Tier        string `json:"tier" validate:"required"`
}
