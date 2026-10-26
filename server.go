package main

import (
	"account-api/graph"
	"account-api/internal/config"
	"account-api/internal/repository"
	"account-api/internal/service"
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const defaultPort = "8081"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := config.NewMySQLConnection(
		getEnv("DB_USER", "root"),
		getEnv("DB_PASSWORD", "root"),
		getEnv("DB_HOST", "127.0.0.1"),
		getEnv("DB_PORT", "3307"),
		getEnv("DB_NAME", "account_db"),
	)
	if err != nil {
		log.Fatal("gagal konek ke database:", err)
	}
	defer db.Close()

	repo := repository.NewAccountRepository(db)
	svc := service.NewAccountService(repo)

	txRepo := repository.NewTransactionRepository(db)
	txSvc := service.NewTransactionService(txRepo)

	resolver := &graph.Resolver{
		AccountService:     svc,
		TransactionService: txSvc,
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Satu tempat untuk menerjemahkan error domain jadi pesan aman + extensions.code.
	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		gqlErr := graphql.DefaultErrorPresenter(ctx, err)

		switch {
		case errors.Is(err, service.ErrInvalidEmail):
			gqlErr.Extensions = map[string]any{"code": "BAD_USER_INPUT"}
		case errors.Is(err, repository.ErrDuplicate):
			gqlErr.Message = "akun sudah ada"
			gqlErr.Extensions = map[string]any{"code": "CONFLICT"}
		case errors.Is(err, repository.ErrNotFound):
			gqlErr.Message = "data tidak ditemukan"
			gqlErr.Extensions = map[string]any{"code": "NOT_FOUND"}
		default:
			// Error bawaan gqlgen (parse/validasi query) sudah punya code sendiri.
			if _, ok := gqlErr.Extensions["code"]; !ok {
				log.Println("internal error:", err)
				gqlErr.Message = "terjadi kesalahan"
				gqlErr.Extensions = map[string]any{"code": "INTERNAL"}
			}
		}

		return gqlErr
	})

	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// http.Handle("/query", srv)

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	// bungkus dengan Middleware supaya SETIAP request punya loader sendiri
	loaderMiddleware := graph.Middleware(txSvc, mux)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, loaderMiddleware))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
