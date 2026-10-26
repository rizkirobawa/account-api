package graph

import "account-api/internal/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AccountService     *service.AccountService
	TransactionService *service.TransactionService
}
