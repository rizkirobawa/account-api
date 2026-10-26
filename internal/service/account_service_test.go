package service

import (
	"context"
	"errors"
	"testing"
)

// Validasi email terjadi sebelum repo dipanggil, jadi repo boleh nil.
func TestCreateAccount_EmptyEmail(t *testing.T) {
	svc := NewAccountService(nil)

	_, err := svc.CreateAccount(context.Background(), "acc-1", "")
	if !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf("err = %v, want ErrInvalidEmail", err)
	}
}
