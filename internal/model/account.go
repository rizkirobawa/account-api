package model

type Account struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Balance int64  `json:"balance"`
}
