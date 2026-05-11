package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Josey34/goshop/delivery/http/middleware"
	"github.com/Josey34/goshop/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func newJWT() *jwt.JWTService {
	return jwt.NewJWTService("test-secret", 1)
}

func setupProtectedRouter(jwtSvc *jwt.JWTService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", middleware.AuthMiddleware(jwtSvc), func(c *gin.Context) {
		id, _ := c.Get("customer_id")
		c.JSON(http.StatusOK, gin.H{"customer_id": id})
	})
	return r
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	svc := newJWT()
	token, err := svc.Generate("c1", "john@example.com")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	r := setupProtectedRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got status=%d, want=200", w.Code)
	}
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	r := setupProtectedRouter(newJWT())
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got status=%d, want=401", w.Code)
	}
}

func TestAuthMiddleware_WrongPrefix(t *testing.T) {
	svc := newJWT()
	token, _ := svc.Generate("c1", "john@example.com")

	r := setupProtectedRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got status=%d, want=401", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	r := setupProtectedRouter(newJWT())
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got status=%d, want=401", w.Code)
	}
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	attacker := jwt.NewJWTService("wrong-secret", 1)
	token, _ := attacker.Generate("c1", "john@example.com")

	r := setupProtectedRouter(newJWT())
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got status=%d, want=401", w.Code)
	}
}

func TestAuthMiddleware_SetsCustomerIDInContext(t *testing.T) {
	svc := newJWT()
	token, _ := svc.Generate("c42", "john@example.com")

	r := setupProtectedRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("got status=%d, want=200", w.Code)
	}
	body := w.Body.String()
	if body == "" {
		t.Error("empty response body")
	}
	// verify customer_id=c42 in response
	if !strings.Contains(body, "c42") {
		t.Errorf("customer_id not set in context, body=%s", body)
	}
}
