package errors_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/errors"
)

func TestDomainError_Error(t *testing.T) {
	t.Run("without fields", func(t *testing.T) {
		err := errors.NewNotFound("product", "p1")
		if err.Error() == "" {
			t.Error("expected non-empty error string")
		}
	})

	t.Run("with fields", func(t *testing.T) {
		err := errors.NewValidation("name", map[string]string{"name": "required"})
		if err.Error() == "" {
			t.Error("expected non-empty error string")
		}
	})
}

func TestIsNotFound(t *testing.T) {
	if !errors.IsNotFound(errors.NewNotFound("product", "p1")) {
		t.Error("expected IsNotFound=true")
	}
	if errors.IsNotFound(errors.NewValidation("x", nil)) {
		t.Error("expected IsNotFound=false for validation error")
	}
}

func TestIsValidation(t *testing.T) {
	if !errors.IsValidation(errors.NewValidation("x", nil)) {
		t.Error("expected IsValidation=true")
	}
	if errors.IsValidation(errors.NewNotFound("product", "p1")) {
		t.Error("expected IsValidation=false for not found error")
	}
}

func TestIsConflict(t *testing.T) {
	if !errors.IsConflict(errors.NewConflict("customer", "john@example.com")) {
		t.Error("expected IsConflict=true")
	}
}

func TestIsUnauthorized(t *testing.T) {
	if !errors.IsUnauthorized(errors.NewUnauthorized("invalid credentials")) {
		t.Error("expected IsUnauthorized=true")
	}
}

func TestIsInsufficientStock(t *testing.T) {
	if !errors.IsInsufficientStock(errors.NewInsufficientStock("p1", 5, 2)) {
		t.Error("expected IsInsufficientStock=true")
	}
}

func TestIsInvalidTransition(t *testing.T) {
	if !errors.IsInvalidTransition(errors.NewInvalidTransition("PENDING", "SHIPPED")) {
		t.Error("expected IsInvalidTransition=true")
	}
}
