package graph

import (
	"account-api/graph/model"
	"account-api/internal/service"
	"context"
	"net/http"

	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

// key untuk nyimpen loader di context request
const loadersKey = "dataloaders"

// Loaders adalah kumpulan semua loader yang tersedia dalam satu request
type Loaders struct {
	TransactionsByAccountID *dataloadgen.Loader[string, []*model.Transaction]
}

func NewLoaders(txService *service.TransactionService) *Loaders {
	return &Loaders{
		TransactionsByAccountID: dataloadgen.NewLoader(
			batchTransactionsByAccoundtID(txService),
			dataloadgen.WithWait(1), // window batching: 1ms
		),
	}
}

// batchTransactionsByAccountID adalah fungsi batch yang dipanggil dataloadgen
// dengan SEMUA account_id yang diminta dalam satu window waktu.
func batchTransactionsByAccoundtID(txService *service.TransactionService) func(ctx context.Context, accountIDs []string) ([][]*model.Transaction, []error) {
	return func(ctx context.Context, accountIDs []string) ([][]*model.Transaction, []error) {
		// satu query untuk SEMUA account_id sekaligus
		transactions, err := txService.ListAccountByIDs(ctx, accountIDs)
		if err != nil {
			errs := make([]error, len(accountIDs))
			for i := range errs {
				errs[i] = err
			}
			return nil, errs
		}

		grouped := make(map[string][]*model.Transaction)
		for _, tx := range transactions {
			tx := tx //hindari gotcha loop variable
			grouped[tx.AccountID] = append(grouped[tx.AccountID], &model.Transaction{
				ID:     tx.ID,
				Amount: int(tx.Amount),
				Type:   tx.Type,
			})
		}

		// PENTING: urutan hasil HARUS sama persis dengan urutan accountIDs yang diminta
		result := make([][]*model.Transaction, len(accountIDs))
		for i, id := range accountIDs {
			result[i] = grouped[id] // kalau nggak ada, otomatis nil (slice kosong)
		}

		return result, nil
	}
}

// Middleware menyisipkan Loaders baru ke context di SETIAP request HTTP masuk
func Middleware(txService *service.TransactionService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loaders := NewLoaders(txService)
		ctx := context.WithValue(r.Context(), loadersKey, loaders)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
