package main

import (
	"account-api/internal/config"
	"account-api/internal/handler"
	"account-api/internal/repository"
	"account-api/internal/service"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := config.NewMySQLConnection(
		getEnv("DB_USER", "root"),
		getEnv("DB_PASSWORD", "123456"),
		getEnv("DB_HOST", "127.0.0.1"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_NAME", "account_db"),
	)
	if err != nil {
		log.Fatal("gagal konek ke database:", err)
	}
	defer db.Close()

	repo := repository.NewAccountRepository(db)
	svc := service.NewAccountService(repo)
	h := handler.NewAccountHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /accounts", h.Create)
	mux.HandleFunc("GET /accounts/{id}", h.Get)
	mux.HandleFunc("GET /accounts", h.List)

	log.Println("Server jalan di port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
