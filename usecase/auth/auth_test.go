package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	ucauth "github.com/Josey34/goshop/usecase/auth"
)

// --- mocks ---

type mockCustomerRepo struct {
	customer      *entity.Customer
	existsByEmail bool
	findErr       error
	existsErr     error
	createErr     error
}

func (m *mockCustomerRepo) Create(ctx context.Context, c *entity.Customer) error {
	return m.createErr
}
func (m *mockCustomerRepo) FindByID(ctx context.Context, id string) (*entity.Customer, error) {
	return nil, nil
}
func (m *mockCustomerRepo) FindByEmail(ctx context.Context, email string) (*entity.Customer, error) {
	return m.customer, m.findErr
}
func (m *mockCustomerRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return m.existsByEmail, m.existsErr
}

type mockHasher struct {
	hashResult string
	hashErr    error
	compareErr error
}

func (m *mockHasher) Hash(password string) (string, error) {
	return m.hashResult, m.hashErr
}
func (m *mockHasher) Compare(hashed, plain string) error {
	return m.compareErr
}

type mockIDGen struct{ id string }

func (m *mockIDGen) Generate() string { return m.id }

func makeCustomer(email, hash string) *entity.Customer {
	e, _ := valueobject.NewEmail(email)
	p, _ := valueobject.NewPhone("08123456789")
	a, _ := valueobject.NewAddress("Jl. A", "Jakarta", "DKI", "10220")
	return &entity.Customer{
		ID:           "c1",
		Name:         "John",
		Email:        e,
		Phone:        p,
		Address:      a,
		PasswordHash: hash,
	}
}

// --- login tests ---

func TestLoginUseCase_Execute(t *testing.T) {
	t.Run("valid credentials returns customer", func(t *testing.T) {
		repo := &mockCustomerRepo{customer: makeCustomer("john@example.com", "hashed")}
		h := &mockHasher{compareErr: nil}
		uc := ucauth.NewLoginUseCase(repo, h)

		customer, err := uc.Execute(context.Background(), ucauth.LoginInput{
			Email: "john@example.com", Password: "secret",
		})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if customer.ID != "c1" {
			t.Errorf("got id=%s, want=c1", customer.ID)
		}
	})

	t.Run("wrong password returns unauthorized", func(t *testing.T) {
		repo := &mockCustomerRepo{customer: makeCustomer("john@example.com", "hashed")}
		h := &mockHasher{compareErr: errors.New("mismatch")}
		uc := ucauth.NewLoginUseCase(repo, h)

		_, err := uc.Execute(context.Background(), ucauth.LoginInput{
			Email: "john@example.com", Password: "wrong",
		})
		if err == nil {
			t.Error("expected error for wrong password")
		}
	})

	t.Run("email not found returns error", func(t *testing.T) {
		repo := &mockCustomerRepo{findErr: errors.New("not found")}
		uc := ucauth.NewLoginUseCase(repo, &mockHasher{})

		_, err := uc.Execute(context.Background(), ucauth.LoginInput{
			Email: "unknown@example.com", Password: "secret",
		})
		if err == nil {
			t.Error("expected error for unknown email")
		}
	})
}

// --- register tests ---

func TestRegisterUseCase_Execute(t *testing.T) {
	validInput := ucauth.RegisterInput{
		Name:       "John",
		Email:      "john@example.com",
		Phone:      "08123456789",
		Street:     "Jl. A",
		City:       "Jakarta",
		Province:   "DKI",
		PostalCode: "10220",
		Password:   "secret",
	}

	t.Run("valid registration creates customer", func(t *testing.T) {
		repo := &mockCustomerRepo{existsByEmail: false}
		h := &mockHasher{hashResult: "hashed_secret"}
		idGen := &mockIDGen{id: "c1"}
		uc := ucauth.NewRegisterUseCase(repo, h, idGen)

		err := uc.Execute(context.Background(), validInput)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
	})

	t.Run("duplicate email returns conflict", func(t *testing.T) {
		repo := &mockCustomerRepo{existsByEmail: true}
		uc := ucauth.NewRegisterUseCase(repo, &mockHasher{hashResult: "h"}, &mockIDGen{id: "c1"})

		err := uc.Execute(context.Background(), validInput)
		if err == nil {
			t.Error("expected conflict error for duplicate email")
		}
	})

	t.Run("invalid email returns error", func(t *testing.T) {
		input := validInput
		input.Email = "not-an-email"
		repo := &mockCustomerRepo{existsByEmail: false}
		uc := ucauth.NewRegisterUseCase(repo, &mockHasher{hashResult: "h"}, &mockIDGen{id: "c1"})

		err := uc.Execute(context.Background(), input)
		if err == nil {
			t.Error("expected error for invalid email")
		}
	})
}
