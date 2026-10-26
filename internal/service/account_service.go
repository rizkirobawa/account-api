package service

import (
	"account-api/internal/model"
	"account-api/internal/repository"
	"context"
	"errors"
)

var ErrInvalidEmail = errors.New("Email tidak valid")

type AccountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, id, email string) (model.Account, error) {
	if email == "" {
		return model.Account{}, ErrInvalidEmail
	}

	acc := model.Account{
		ID:      id,
		Email:   email,
		Balance: 0,
	}

	if err := s.repo.Create(ctx, acc); err != nil {
		return model.Account{}, err
	}
	return acc, nil
}

func (s *AccountService) GetAccount(ctx context.Context, id string) (model.Account, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *AccountService) ListAccounts(ctx context.Context) ([]model.Account, error) {
	return s.repo.FindByAll(ctx)
}
