package jwt_test

import (
	"testing"

	"github.com/Josey34/goshop/pkg/jwt"
)

func TestJWTService_GenerateAndValidate(t *testing.T) {
	svc := jwt.NewJWTService("test-secret", 24)

	t.Run("generate returns non-empty token", func(t *testing.T) {
		token, err := svc.Generate("c1", "john@example.com")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if token == "" {
			t.Error("expected non-empty token")
		}
	})

	t.Run("validate returns correct claims", func(t *testing.T) {
		token, _ := svc.Generate("c1", "john@example.com")
		claims, err := svc.Validate(token)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if claims.CustomerID != "c1" {
			t.Errorf("got customerID=%s, want=c1", claims.CustomerID)
		}
		if claims.Email != "john@example.com" {
			t.Errorf("got email=%s, want=john@example.com", claims.Email)
		}
	})

	t.Run("validate with wrong secret returns error", func(t *testing.T) {
		token, _ := svc.Generate("c1", "john@example.com")
		otherSvc := jwt.NewJWTService("wrong-secret", 24)
		_, err := otherSvc.Validate(token)
		if err == nil {
			t.Error("expected error for wrong secret")
		}
	})

	t.Run("validate with tampered token returns error", func(t *testing.T) {
		_, err := svc.Validate("invalid.token.string")
		if err == nil {
			t.Error("expected error for invalid token")
		}
	})
}
