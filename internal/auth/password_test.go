package auth

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hashedPassword == "" {
		t.Error("HashPassword returned empty hash")
	}

	if hashedPassword == password {
		t.Error("HashPassword should not return the original password")
	}

	if !strings.HasPrefix(hashedPassword, "$2a$") {
		t.Error("HashPassword should return a bcrypt hash starting with $2a$")
	}
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	password := ""

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed with empty password: %v", err)
	}

	if hashedPassword == "" {
		t.Error("HashPassword returned empty hash for empty password")
	}
}

func TestHashPassword_DifferentHashes(t *testing.T) {
	password := "testpassword123"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash1 == hash2 {
		t.Error("HashPassword should generate different hashes for the same password (due to salt)")
	}
}

func TestComparePassword_ValidPassword(t *testing.T) {
	password := "testpassword123"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	err = ComparePassword(hashedPassword, password)
	if err != nil {
		t.Errorf("ComparePassword failed with valid password: %v", err)
	}
}

func TestComparePassword_InvalidPassword(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	err = ComparePassword(hashedPassword, wrongPassword)
	if err == nil {
		t.Error("ComparePassword should have failed with wrong password")
	}
}

func TestComparePassword_EmptyPassword(t *testing.T) {
	password := "testpassword123"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	err = ComparePassword(hashedPassword, "")
	if err == nil {
		t.Error("ComparePassword should have failed with empty password")
	}
}

func TestComparePassword_InvalidHash(t *testing.T) {
	password := "testpassword123"
	invalidHash := "invalid-hash"

	err := ComparePassword(invalidHash, password)
	if err == nil {
		t.Error("ComparePassword should have failed with invalid hash")
	}
}

func TestHashAndComparePassword_RoundTrip(t *testing.T) {
	testCases := []string{
		"password123",
		"ComplexP@ssw0rd!",
		"simple",
		"verylongpasswordwithmanyCHARACTERSandNumbers123456789",
		"パスワード", // Unicode characters
		"",          // Empty password
	}

	for _, password := range testCases {
		t.Run("password: "+password, func(t *testing.T) {
			hashedPassword, err := HashPassword(password)
			if err != nil {
				t.Fatalf("HashPassword failed: %v", err)
			}

			err = ComparePassword(hashedPassword, password)
			if err != nil {
				t.Errorf("ComparePassword failed with original password: %v", err)
			}
		})
	}
}