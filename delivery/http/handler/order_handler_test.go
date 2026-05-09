package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Josey34/goshop/delivery/http/handler"
	"github.com/Josey34/goshop/delivery/http/middleware"
	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	ucorder "github.com/Josey34/goshop/usecase/order"
	"github.com/gin-gonic/gin"
)

type mockOrderService struct {
	createFn  func(ctx context.Context, input ucorder.CreateOrderInput) (*entity.Order, error)
	getByIDFn func(ctx context.Context, id string) (*entity.Order, error)
	listFn    func(ctx context.Context, customerID string, pagination valueobject.Pagination) ([]*entity.Order, error)
}

func (m *mockOrderService) CreateOrder(ctx context.Context, input ucorder.CreateOrderInput) (*entity.Order, error) {
	return m.createFn(ctx, input)
}
func (m *mockOrderService) GetByID(ctx context.Context, id string) (*entity.Order, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockOrderService) ListByCustomer(ctx context.Context, customerID string, pagination valueobject.Pagination) ([]*entity.Order, error) {
	return m.listFn(ctx, customerID, pagination)
}

func makeOrder(id, customerID string) *entity.Order {
	return &entity.Order{
		ID:         id,
		CustomerID: customerID,
		Status:     valueobject.OrderStatusPending,
		Total:      valueobject.NewMoney(2000),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func setupOrderRouter(h *handler.OrderHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorMiddleware())
	r.POST("/orders", func(c *gin.Context) {
		c.Set("customer_id", "c1")
		h.Create(c)
	})
	r.GET("/orders/:id", h.GetByID)
	r.GET("/orders", func(c *gin.Context) {
		c.Set("customer_id", "c1")
		h.List(c)
	})
	return r
}

func TestOrderHandler_Create(t *testing.T) {
	t.Run("201 on valid order", func(t *testing.T) {
		svc := &mockOrderService{
			createFn: func(ctx context.Context, input ucorder.CreateOrderInput) (*entity.Order, error) {
				return makeOrder("o1", "c1"), nil
			},
		}
		h := handler.NewOrderHandler(svc)
		r := setupOrderRouter(h)

		body, _ := json.Marshal(map[string]any{
			"items": []map[string]any{{"product_id": "p1", "quantity": 2}},
		})
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("got status=%d, want=201", w.Code)
		}
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["id"] != "o1" {
			t.Errorf("got id=%v, want=o1", resp["id"])
		}
	})

	t.Run("400 on empty items", func(t *testing.T) {
		svc := &mockOrderService{}
		h := handler.NewOrderHandler(svc)
		r := setupOrderRouter(h)

		body, _ := json.Marshal(map[string]any{"items": []any{}})
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code == http.StatusCreated {
			t.Error("expected non-201 for empty items")
		}
	})
}

func TestOrderHandler_GetByID(t *testing.T) {
	t.Run("200 with order", func(t *testing.T) {
		svc := &mockOrderService{
			getByIDFn: func(ctx context.Context, id string) (*entity.Order, error) {
				return makeOrder(id, "c1"), nil
			},
		}
		h := handler.NewOrderHandler(svc)
		r := setupOrderRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/orders/o1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status=%d, want=200", w.Code)
		}
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["id"] != "o1" {
			t.Errorf("got id=%v, want=o1", resp["id"])
		}
	})
}

func TestOrderHandler_List(t *testing.T) {
	t.Run("200 with orders list", func(t *testing.T) {
		svc := &mockOrderService{
			listFn: func(ctx context.Context, customerID string, pagination valueobject.Pagination) ([]*entity.Order, error) {
				return []*entity.Order{makeOrder("o1", customerID), makeOrder("o2", customerID)}, nil
			},
		}
		h := handler.NewOrderHandler(svc)
		r := setupOrderRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/orders", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status=%d, want=200", w.Code)
		}
		var resp []any
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp) != 2 {
			t.Errorf("got %d orders, want=2", len(resp))
		}
	})
}
