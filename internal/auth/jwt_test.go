package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sanntintdev/chirpy/internal/database"
)

func TestNewTokenService(t *testing.T) {
	secretKey := []byte("test-secret-key")
	issuer := "test-issuer"

	ts := NewTokenService(secretKey, issuer)

	if ts == nil {
		t.Fatal("NewTokenService returned nil")
	}
	if string(ts.secretKey) != string(secretKey) {
		t.Errorf("Expected secret key %s, got %s", secretKey, ts.secretKey)
	}
	if ts.issuer != issuer {
		t.Errorf("Expected issuer %s, got %s", issuer, ts.issuer)
	}
}

func TestGenerateToken(t *testing.T) {
	secretKey := []byte("test-secret-key")
	issuer := "test-issuer"
	ts := NewTokenService(secretKey, issuer)

	userID := uuid.New()
	user := &database.User{
		ID: userID,
	}
	expiration := time.Hour

	token, err := ts.GenerateToken(user, expiration)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("GenerateToken returned empty token")
	}
}

func TestValidateToken_ValidToken(t *testing.T) {
	secretKey := []byte("test-secret-key")
	issuer := "test-issuer"
	ts := NewTokenService(secretKey, issuer)

	userID := uuid.New()
	user := &database.User{
		ID: userID,
	}
	expiration := time.Hour

	token, err := ts.GenerateToken(user, expiration)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	validatedUserID, err := ts.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, validatedUserID)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	secretKey := []byte("test-secret-key")
	issuer := "test-issuer"
	ts := NewTokenService(secretKey, issuer)

	invalidToken := "invalid.token.string"

	_, err := ts.ValidateToken(invalidToken)
	if err == nil {
		t.Error("ValidateToken should have failed with invalid token")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	secretKey := []byte("test-secret-key")
	issuer := "test-issuer"
	ts := NewTokenService(secretKey, issuer)

	userID := uuid.New()
	user := &database.User{
		ID: userID,
	}
	expiration := -time.Hour // Expired token

	token, err := ts.GenerateToken(user, expiration)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = ts.ValidateToken(token)
	if err == nil {
		t.Error("ValidateToken should have failed with expired token")
	}
}

func TestValidateToken_WrongSecretKey(t *testing.T) {
	secretKey := []byte("test-secret-key")
	wrongSecretKey := []byte("wrong-secret-key")
	issuer := "test-issuer"
	
	ts1 := NewTokenService(secretKey, issuer)
	ts2 := NewTokenService(wrongSecretKey, issuer)

	userID := uuid.New()
	user := &database.User{
		ID: userID,
	}
	expiration := time.Hour

	token, err := ts1.GenerateToken(user, expiration)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = ts2.ValidateToken(token)
	if err == nil {
		t.Error("ValidateToken should have failed with wrong secret key")
	}
}