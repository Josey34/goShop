package factory_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/factory"
)

func TestNewCustomer(t *testing.T) {
	validInput := factory.CreateCustomerInput{
		Name:       "John Doe",
		Email:      "john@example.com",
		Phone:      "08123456789",
		Street:     "Jl. Sudirman",
		City:       "Jakarta",
		Province:   "DKI Jakarta",
		PostalCode: "10220",
		Password:   "secret",
	}

	t.Run("valid customer", func(t *testing.T) {
		c, err := factory.NewCustomer("c1", validInput, "hashed_secret")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if c.ID != "c1" {
			t.Errorf("got id=%s, want=c1", c.ID)
		}
		if c.PasswordHash != "hashed_secret" {
			t.Errorf("got hash=%s, want=hashed_secret", c.PasswordHash)
		}
		if c.Email.Value() != "john@example.com" {
			t.Errorf("got email=%s, want=john@example.com", c.Email.Value())
		}
	})

	t.Run("invalid email returns error", func(t *testing.T) {
		input := validInput
		input.Email = "not-an-email"
		_, err := factory.NewCustomer("c1", input, "hash")
		if err == nil {
			t.Error("expected error for invalid email")
		}
	})

	t.Run("invalid phone returns error", func(t *testing.T) {
		input := validInput
		input.Phone = "12345"
		_, err := factory.NewCustomer("c1", input, "hash")
		if err == nil {
			t.Error("expected error for invalid phone")
		}
	})

	t.Run("invalid address returns error", func(t *testing.T) {
		input := validInput
		input.PostalCode = "123"
		_, err := factory.NewCustomer("c1", input, "hash")
		if err == nil {
			t.Error("expected error for invalid postal code")
		}
	})
}
