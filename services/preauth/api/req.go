package api

type PreauthReq struct {
	AccountID string `json:"account_id" validate:"required"`
}
