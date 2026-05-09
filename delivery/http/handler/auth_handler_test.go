package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Josey34/goshop/delivery/http/handler"
	"github.com/Josey34/goshop/delivery/http/middleware"
	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	"github.com/Josey34/goshop/service"
	ucauth "github.com/Josey34/goshop/usecase/auth"
	"github.com/gin-gonic/gin"
)

type mockAuthService struct {
	registerFn func(ctx context.Context, input ucauth.RegisterInput) error
	loginFn    func(ctx context.Context, input ucauth.LoginInput) (*service.LoginOutput, error)
}

func (m *mockAuthService) Register(ctx context.Context, input ucauth.RegisterInput) error {
	return m.registerFn(ctx, input)
}
func (m *mockAuthService) Login(ctx context.Context, input ucauth.LoginInput) (*service.LoginOutput, error) {
	return m.loginFn(ctx, input)
}

func setupAuthRouter(h *handler.AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorMiddleware())
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	return r
}

func makeCustomer(id, email string) *entity.Customer {
	e, _ := valueobject.NewEmail(email)
	return &entity.Customer{ID: id, Email: e}
}

func TestAuthHandler_Register(t *testing.T) {
	validBody := map[string]any{
		"name": "John", "email": "john@example.com",
		"password": "secret123", "phone": "08123456789",
	}

	t.Run("201 on valid registration", func(t *testing.T) {
		svc := &mockAuthService{
			registerFn: func(ctx context.Context, input ucauth.RegisterInput) error { return nil },
		}
		h := handler.NewAuthHandler(svc)
		r := setupAuthRouter(h)

		body, _ := json.Marshal(validBody)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("got status=%d, want=201", w.Code)
		}
	})

	t.Run("400 on missing fields", func(t *testing.T) {
		svc := &mockAuthService{}
		h := handler.NewAuthHandler(svc)
		r := setupAuthRouter(h)

		body, _ := json.Marshal(map[string]any{"name": "John"})
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code == http.StatusCreated {
			t.Error("expected non-201 for missing fields")
		}
	})

	t.Run("service error propagates", func(t *testing.T) {
		svc := &mockAuthService{
			registerFn: func(ctx context.Context, input ucauth.RegisterInput) error {
				return errors.New("conflict")
			},
		}
		h := handler.NewAuthHandler(svc)
		r := setupAuthRouter(h)

		body, _ := json.Marshal(validBody)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code == http.StatusCreated {
			t.Error("expected non-201 for service error")
		}
	})
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("200 with token on valid login", func(t *testing.T) {
		svc := &mockAuthService{
			loginFn: func(ctx context.Context, input ucauth.LoginInput) (*service.LoginOutput, error) {
				return &service.LoginOutput{
					Token:    "jwt-token",
					Customer: makeCustomer("c1", "john@example.com"),
				}, nil
			},
		}
		h := handler.NewAuthHandler(svc)
		r := setupAuthRouter(h)

		body, _ := json.Marshal(map[string]any{"email": "john@example.com", "password": "secret123"})
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status=%d, want=200", w.Code)
		}
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["token"] != "jwt-token" {
			t.Errorf("got token=%v, want=jwt-token", resp["token"])
		}
	})

	t.Run("400 on missing password", func(t *testing.T) {
		svc := &mockAuthService{}
		h := handler.NewAuthHandler(svc)
		r := setupAuthRouter(h)

		body, _ := json.Marshal(map[string]any{"email": "john@example.com"})
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Error("expected non-200 for missing password")
		}
	})
}
