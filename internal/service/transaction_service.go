package service

import (
	"account-api/internal/model"
	"account-api/internal/repository"
	"context"
)

type TransactionService struct {
	repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) ListByAccount(ctx context.Context, accountID string) ([]model.Transaction, error) {
	return s.repo.FindByAccountId(ctx, accountID)
}

func (s *TransactionService) ListAccountByIDs(ctx context.Context, accountIDs []string) ([]model.Transaction, error) {
	return s.repo.FindByAccountIDs(ctx, accountIDs)
}
