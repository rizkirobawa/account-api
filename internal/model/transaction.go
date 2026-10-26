package model

type Transaction struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	Amount    int64  `json:"amount"`
	Type      string `json:"type"`
}
