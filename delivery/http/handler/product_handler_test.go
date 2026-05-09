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
	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	ucproduct "github.com/Josey34/goshop/usecase/product"
	"github.com/gin-gonic/gin"
)

// --- mock service ---

type mockProductService struct {
	createFn  func(ctx context.Context, input ucproduct.CreateProductInput) (*entity.Product, error)
	getByIDFn func(ctx context.Context, id string) (*entity.Product, error)
}

func (m *mockProductService) Create(ctx context.Context, input ucproduct.CreateProductInput) (*entity.Product, error) {
	return m.createFn(ctx, input)
}
func (m *mockProductService) GetByID(ctx context.Context, id string) (*entity.Product, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockProductService) List(ctx context.Context, pagination valueobject.Pagination) ([]*entity.Product, error) {
	return nil, nil
}
func (m *mockProductService) Update(ctx context.Context, input ucproduct.UpdateProductInput) (*entity.Product, error) {
	return nil, nil
}
func (m *mockProductService) Delete(ctx context.Context, id string) error { return nil }
func (m *mockProductService) UploadImage(ctx context.Context, image ucproduct.UploadProductImage) (string, error) {
	return "", nil
}

// --- helpers ---

func setupRouter(h *handler.ProductHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/products", h.Create)
	r.GET("/products/:id", h.GetByID)
	return r
}

func makeProduct(id, name string, price int64) *entity.Product {
	return &entity.Product{
		ID:        id,
		Name:      name,
		Price:     valueobject.NewMoney(price),
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// --- tests ---

func TestProductHandler_Create(t *testing.T) {
	t.Run("201 on valid product", func(t *testing.T) {
		svc := &mockProductService{
			createFn: func(ctx context.Context, input ucproduct.CreateProductInput) (*entity.Product, error) {
				return makeProduct("p1", input.Name, input.Price), nil
			},
		}
		h := handler.NewProductHandler(svc)
		r := setupRouter(h)

		body, _ := json.Marshal(map[string]any{"name": "Shirt", "price": 1000, "stock": 5})
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("got status=%d, want=201", w.Code)
		}

		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["name"] != "Shirt" {
			t.Errorf("got name=%v, want=Shirt", resp["name"])
		}
	})

	t.Run("400 on missing body", func(t *testing.T) {
		svc := &mockProductService{}
		h := handler.NewProductHandler(svc)
		r := setupRouter(h)

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code == http.StatusCreated {
			t.Error("expected non-201 for bad body")
		}
	})
}

func TestProductHandler_GetByID(t *testing.T) {
	t.Run("200 with product", func(t *testing.T) {
		svc := &mockProductService{
			getByIDFn: func(ctx context.Context, id string) (*entity.Product, error) {
				return makeProduct(id, "Shirt", 1000), nil
			},
		}
		h := handler.NewProductHandler(svc)
		r := setupRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/products/p1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status=%d, want=200", w.Code)
		}

		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["id"] != "p1" {
			t.Errorf("got id=%v, want=p1", resp["id"])
		}
	})
}
