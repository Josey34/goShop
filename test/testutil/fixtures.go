package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func RegisterAndLogin(t *testing.T, engine *gin.Engine) string {
	t.Helper()

	regBody, _ := json.Marshal(map[string]any{
		"name": "Test User", "email": "test@example.com",
		"password": "secret123", "phone": "08123456789",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("register failed: status=%d body=%s", w.Code, w.Body.String())
	}

	loginBody, _ := json.Marshal(map[string]any{
		"email": "test@example.com", "password": "secret123",
	})
	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login failed: status=%d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	token, ok := resp["token"].(string)
	if !ok || token == "" {
		t.Fatalf("no token in login response: %s", w.Body.String())
	}
	return token
}

func CreateProduct(t *testing.T, engine *gin.Engine, token string, name string, price int64, stock int) string {
	t.Helper()

	body, _ := json.Marshal(map[string]any{
		"name": name, "price": price, "stock": stock, "description": "test product",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create product failed: status=%d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	id, ok := resp["id"].(string)
	if !ok || id == "" {
		t.Fatalf("no id in create product response: %s", w.Body.String())
	}
	return id
}
