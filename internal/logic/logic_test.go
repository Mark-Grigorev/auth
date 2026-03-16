package logic

import (
	"context"
	"errors"
	"testing"

	"github.com/Mark-Grigorev/auth/internal/model"
)

// --- mocks ---

type mockDB struct {
	createUserFn   func(ctx context.Context, userData *model.UserRegistrationData) (int64, error)
	authorisationFn func(ctx context.Context, login, password string) (int64, error)
}

func (m *mockDB) CreateUser(ctx context.Context, userData *model.UserRegistrationData) (int64, error) {
	return m.createUserFn(ctx, userData)
}

func (m *mockDB) Authorisation(ctx context.Context, login, password string) (int64, error) {
	return m.authorisationFn(ctx, login, password)
}

type mockJWT struct {
	createTokenFn   func(userID int64) (string, error)
	validateTokenFn func(token string) (bool, error)
}

func (m *mockJWT) CreateToken(userID int64) (string, error) {
	return m.createTokenFn(userID)
}

func (m *mockJWT) ValidateToken(token string) (bool, error) {
	return m.validateTokenFn(token)
}

// --- helpers ---

func newLogic(db DBProvider, jwt JWTProvider) *Logic {
	return New(&model.Config{}, db, jwt)
}

// --- Register ---

func TestRegister_Success(t *testing.T) {
	db := &mockDB{
		createUserFn: func(_ context.Context, _ *model.UserRegistrationData) (int64, error) {
			return 42, nil
		},
	}
	l := newLogic(db, &mockJWT{})

	id, err := l.Register(context.Background(), &model.UserRegistrationData{
		Login:    "user",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 42 {
		t.Fatalf("expected id 42, got %d", id)
	}
}

func TestRegister_DBError(t *testing.T) {
	db := &mockDB{
		createUserFn: func(_ context.Context, _ *model.UserRegistrationData) (int64, error) {
			return 0, errors.New("db error")
		},
	}
	l := newLogic(db, &mockJWT{})

	_, err := l.Register(context.Background(), &model.UserRegistrationData{Password: "p"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRegister_PasswordIsHashed(t *testing.T) {
	var savedPass string
	db := &mockDB{
		createUserFn: func(_ context.Context, u *model.UserRegistrationData) (int64, error) {
			savedPass = u.Password
			return 1, nil
		},
	}
	l := newLogic(db, &mockJWT{})

	_, err := l.Register(context.Background(), &model.UserRegistrationData{Password: "plaintext"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if savedPass == "plaintext" {
		t.Fatal("password was not hashed before saving")
	}
	if savedPass == "" {
		t.Fatal("saved password is empty")
	}
}

// --- Authorization ---

func TestAuthorization_Success(t *testing.T) {
	db := &mockDB{
		authorisationFn: func(_ context.Context, _, _ string) (int64, error) {
			return 7, nil
		},
	}
	jwt := &mockJWT{
		createTokenFn: func(userID int64) (string, error) {
			if userID != 7 {
				return "", errors.New("wrong userID")
			}
			return "token123", nil
		},
	}
	l := newLogic(db, jwt)

	token, err := l.Authorization(context.Background(), "login", "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "token123" {
		t.Fatalf("expected token123, got %s", token)
	}
}

func TestAuthorization_DBError(t *testing.T) {
	db := &mockDB{
		authorisationFn: func(_ context.Context, _, _ string) (int64, error) {
			return 0, errors.New("not found")
		},
	}
	l := newLogic(db, &mockJWT{})

	_, err := l.Authorization(context.Background(), "login", "pass")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAuthorization_JWTError(t *testing.T) {
	db := &mockDB{
		authorisationFn: func(_ context.Context, _, _ string) (int64, error) {
			return 1, nil
		},
	}
	jwt := &mockJWT{
		createTokenFn: func(_ int64) (string, error) {
			return "", errors.New("jwt error")
		},
	}
	l := newLogic(db, jwt)

	_, err := l.Authorization(context.Background(), "login", "pass")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- ValidateToken ---

func TestValidateToken_Valid(t *testing.T) {
	jwt := &mockJWT{
		validateTokenFn: func(_ string) (bool, error) {
			return true, nil
		},
	}
	l := newLogic(&mockDB{}, jwt)

	ok, err := l.ValidateToken(context.Background(), "sometoken")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected valid=true")
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	jwt := &mockJWT{
		validateTokenFn: func(_ string) (bool, error) {
			return false, errors.New("invalid token")
		},
	}
	l := newLogic(&mockDB{}, jwt)

	ok, err := l.ValidateToken(context.Background(), "badtoken")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if ok {
		t.Fatal("expected valid=false")
	}
}
