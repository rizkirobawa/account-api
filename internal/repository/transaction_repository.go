package repository

import (
	"account-api/internal/model"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// FindByAccountId mengambil semua transaksi milik 1 akun
func (r *TransactionRepository) FindByAccountId(ctx context.Context, accountID string) ([]model.Transaction, error) {
	query := "SELECT id, account_id, amount, type FROM transactions WHERE account_id = ?"

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var tx model.Transaction
		if err := rows.Scan(&tx.ID, &tx.AccountID, &tx.Amount, &tx.Type); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (r *TransactionRepository) FindByAccountIDs(ctx context.Context, accountIDs []string) ([]model.Transaction, error) {
	log.Println("BATCH QUERY untuk account_ids:", accountIDs)
	if len(accountIDs) == 0 {
		return nil, nil
	}

	// bikin placeholder "?, ?, ?" sejumlah accountIDs
	placeholders := make([]string, len(accountIDs))
	args := make([]interface{}, len(accountIDs))
	for i, id := range accountIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT id, account_id, amount, type FROM transactions WHERE account_id IN (%s)`,
		strings.Join(placeholders, ","),
	)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var tx model.Transaction
		if err := rows.Scan(&tx.ID, &tx.AccountID, &tx.Amount, &tx.Type); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
