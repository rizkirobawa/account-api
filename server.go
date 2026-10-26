package main

import (
	"account-api/graph"
	"account-api/internal/config"
	"account-api/internal/repository"
	"account-api/internal/service"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
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
