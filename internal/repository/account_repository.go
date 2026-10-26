package repository

import (
	"account-api/internal/model"
	"context"
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("account not found")

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, acc model.Account) error {
	query := "INSERT INTO accounts (id, email, balance) VALUES (?, ?, ?)"

	_, err := r.db.ExecContext(ctx, query, acc.ID, acc.Email, acc.Balance)
	return err
}

func (r *AccountRepository) FindByID(ctx context.Context, id string) (model.Account, error) {
	query := "SELECT id, email, balance FROM accounts WHERE id = ?"

	var acc model.Account
	err := r.db.QueryRowContext(ctx, query, id).Scan(&acc.ID, &acc.Email, &acc.Balance)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Account{}, ErrNotFound
	}
	if err != nil {
		return model.Account{}, err
	}
	return acc, nil
}

func (r *AccountRepository) FindByAll(ctx context.Context) ([]model.Account, error) {
	query := "SELECT id, email, balance FROM accounts"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // WAJIB, kalau nggak koneksi bocor

	var accounts []model.Account
	for rows.Next() {
		var acc model.Account
		if err := rows.Scan(&acc.ID, &acc.Email, &acc.Balance); err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	// cek error yang mungkin terjadi SELAMA iterasi (bukan cuma saat Query)
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
