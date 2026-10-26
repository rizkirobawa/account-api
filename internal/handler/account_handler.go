package handler

import (
	"account-api/internal/repository"
	"account-api/internal/service"
	"encoding/json"
	"errors"
	"net/http"
)

type AccountHandler struct {
	service *service.AccountService
}

func NewAccountHandler(service *service.AccountService) *AccountHandler {
	return &AccountHandler{
		service: service,
	}
}

type createAccountRequst struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequst
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body tidak valid", http.StatusBadRequest)
		return
	}

	acc, err := h.service.CreateAccount(r.Context(), req.ID, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(acc)
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	acc, err := h.service.GetAccount(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "akun tidak ditemukan", http.StatusNotFound)
			return
		}
		http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.ListAccounts(r.Context())
	if err != nil {
		http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}
