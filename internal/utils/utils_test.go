package utils

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPass_ReturnsHash(t *testing.T) {
	hash, err := HashPass("mypassword")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hash) == 0 {
		t.Fatal("hash is empty")
	}
}

func TestHashPass_IsNotPlaintext(t *testing.T) {
	password := "secret123"
	hash, err := HashPass(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(hash) == password {
		t.Fatal("hash equals plaintext password")
	}
}

func TestHashPass_ValidBcryptHash(t *testing.T) {
	password := "testpass"
	hash, err := HashPass(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		t.Fatalf("hash does not match original password: %v", err)
	}
}

func TestHashPass_DifferentHashesForSamePassword(t *testing.T) {
	h1, _ := HashPass("same")
	h2, _ := HashPass("same")
	if string(h1) == string(h2) {
		t.Fatal("bcrypt should produce different hashes for the same input (random salt)")
	}
}

func TestHashPass_EmptyPassword(t *testing.T) {
	hash, err := HashPass("")
	if err != nil {
		t.Fatalf("unexpected error for empty password: %v", err)
	}
	if len(hash) == 0 {
		t.Fatal("hash is empty for empty password")
	}
}
