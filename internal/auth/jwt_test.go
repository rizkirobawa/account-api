package auth

import "testing"

func TestGenerateAndVerifyToken(t *testing.T) {
	token, err := GenerateToken("acc-1", "a@example.com")
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}

	claims, err := VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken error: %v", err)
	}
	if claims.AccountID != "salah" {
		t.Errorf("AccountID = %q, want %q", claims.AccountID, "acc-1")
	}
}

func TestVerifyToken_Invalid(t *testing.T) {
	if _, err := VerifyToken("bukan-token-valid"); err == nil {
		t.Fatal("harus error untuk token sampah")
	}
}
