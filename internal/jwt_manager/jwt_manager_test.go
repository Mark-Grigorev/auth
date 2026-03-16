package jwtmanager

import (
	"testing"
	"time"
)

func newManager() *Manager {
	return New("test-secret", 5*time.Minute)
}

func TestCreateToken_ReturnsNonEmptyToken(t *testing.T) {
	m := newManager()
	token, err := m.CreateToken(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("token is empty")
	}
}

func TestCreateToken_DifferentUsersGetDifferentTokens(t *testing.T) {
	m := newManager()
	t1, _ := m.CreateToken(1)
	t2, _ := m.CreateToken(2)
	if t1 == t2 {
		t.Fatal("different user IDs produced identical tokens")
	}
}

func TestValidateToken_ValidToken(t *testing.T) {
	m := newManager()
	token, err := m.CreateToken(42)
	if err != nil {
		t.Fatalf("create token error: %v", err)
	}

	ok, err := m.ValidateToken(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected token to be valid")
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	m := newManager()
	ok, err := m.ValidateToken("this.is.not.a.valid.token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
	if ok {
		t.Fatal("expected valid=false for invalid token")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	m1 := New("secret-one", 5*time.Minute)
	m2 := New("secret-two", 5*time.Minute)

	token, err := m1.CreateToken(1)
	if err != nil {
		t.Fatalf("create token error: %v", err)
	}

	ok, err := m2.ValidateToken(token)
	if err == nil {
		t.Fatal("expected error when validating with wrong secret")
	}
	if ok {
		t.Fatal("expected valid=false with wrong secret")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	m := New("secret", -1*time.Second) // уже истёк
	token, err := m.CreateToken(1)
	if err != nil {
		t.Fatalf("create token error: %v", err)
	}

	ok, err := m.ValidateToken(token)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
	if ok {
		t.Fatal("expected valid=false for expired token")
	}
}
