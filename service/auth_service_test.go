package service_test

import (
	"context"
	"testing"

	"github.com/Josey34/goshop/pkg/hasher"
	"github.com/Josey34/goshop/pkg/idgen"
	"github.com/Josey34/goshop/pkg/jwt"
	"github.com/Josey34/goshop/repository/memory"
	"github.com/Josey34/goshop/service"
	ucauth "github.com/Josey34/goshop/usecase/auth"
)

func buildAuthService() *service.AuthService {
	customerRepo := memory.NewCustomerRepo()
	h := hasher.NewBcryptHasher()
	idGen := idgen.NewUUIDGenerator()
	jwtSvc := jwt.NewJWTService("test-secret", 24)

	registerUC := ucauth.NewRegisterUseCase(customerRepo, h, idGen)
	loginUC := ucauth.NewLoginUseCase(customerRepo, h)

	return service.NewAuthService(registerUC, loginUC, jwtSvc)
}

var validRegisterInput = ucauth.RegisterInput{
	Name:       "John",
	Email:      "john@example.com",
	Phone:      "08123456789",
	Street:     "Jl. A",
	City:       "Jakarta",
	Province:   "DKI",
	PostalCode: "10220",
	Password:   "secret123",
}

func TestAuthService_Register(t *testing.T) {
	svc := buildAuthService()

	t.Run("valid registration", func(t *testing.T) {
		if err := svc.Register(context.Background(), validRegisterInput); err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
	})

	t.Run("duplicate email returns conflict", func(t *testing.T) {
		svc := buildAuthService()
		svc.Register(context.Background(), validRegisterInput)
		err := svc.Register(context.Background(), validRegisterInput)
		if err == nil {
			t.Error("expected conflict error")
		}
	})
}

func TestAuthService_Login(t *testing.T) {
	t.Run("valid login returns token", func(t *testing.T) {
		svc := buildAuthService()
		svc.Register(context.Background(), validRegisterInput)

		out, err := svc.Login(context.Background(), ucauth.LoginInput{
			Email:    "john@example.com",
			Password: "secret123",
		})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if out.Token == "" {
			t.Error("expected non-empty token")
		}
		if out.Customer.Email.Value() != "john@example.com" {
			t.Errorf("got email=%s", out.Customer.Email.Value())
		}
	})

	t.Run("wrong password returns error", func(t *testing.T) {
		svc := buildAuthService()
		svc.Register(context.Background(), validRegisterInput)

		_, err := svc.Login(context.Background(), ucauth.LoginInput{
			Email:    "john@example.com",
			Password: "wrongpass",
		})
		if err == nil {
			t.Error("expected error for wrong password")
		}
	})
}
