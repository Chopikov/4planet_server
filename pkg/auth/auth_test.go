package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	service := &Service{}

	password := "testpassword123"
	hash, err := service.HashPassword(password)

	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword returned empty hash")
	}

	if hash == password {
		t.Error("HashPassword returned password as hash")
	}
}

func TestCheckPassword(t *testing.T) {
	service := &Service{}

	password := "testpassword123"
	hash, err := service.HashPassword(password)

	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}

	// Test correct password
	if !service.CheckPassword(password, hash) {
		t.Error("CheckPassword failed for correct password")
	}

	// Test incorrect password
	if service.CheckPassword("wrongpassword", hash) {
		t.Error("CheckPassword succeeded for incorrect password")
	}
}

func TestGenerateToken(t *testing.T) {
	service := &Service{}

	token1 := service.GenerateToken()
	token2 := service.GenerateToken()

	if token1 == "" {
		t.Error("GenerateToken returned empty token")
	}

	if token1 == token2 {
		t.Error("GenerateToken returned duplicate tokens")
	}

	// Check token length (32 bytes = 64 hex chars)
	if len(token1) != 64 {
		t.Errorf("GenerateToken returned token with wrong length: %d", len(token1))
	}
}
