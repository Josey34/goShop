package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Josey34/goshop/test/testutil"
)

// Flow 1: register → login → product CRUD
func TestFlow_ProductCRUD(t *testing.T) {
	engine := testutil.NewTestEngine(t)
	token := testutil.RegisterAndLogin(t, engine)

	// create
	productID := testutil.CreateProduct(t, engine, token, "Shirt", 50000, 10)
	if productID == "" {
		t.Fatal("expected non-empty product ID")
	}

	// list — product appears
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	testutil.AssertStatus(t, w, http.StatusOK)
	list := testutil.ParseJSONArray(t, w)
	if len(list) != 1 {
		t.Errorf("got %d products, want=1", len(list))
	}

	// get by ID
	req = httptest.NewRequest(http.MethodGet, "/api/v1/products/"+productID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	testutil.AssertStatus(t, w, http.StatusOK)
	testutil.AssertJSONField(t, w, "id", productID)

	// update
	updateBody, _ := json.Marshal(map[string]any{"name": "Updated Shirt", "price": 60000, "stock": 5})
	req = httptest.NewRequest(http.MethodPut, "/api/v1/products/"+productID, bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	testutil.AssertStatus(t, w, http.StatusOK)
	testutil.AssertJSONField(t, w, "name", "Updated Shirt")

	// delete
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+productID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	testutil.AssertStatus(t, w, http.StatusNoContent)

	// get after delete → 404
	req = httptest.NewRequest(http.MethodGet, "/api/v1/products/"+productID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	testutil.AssertStatus(t, w, http.StatusNotFound)
}

// Flow 2: register → login → create product → create order → verify total
func TestFlow_OrderTotal(t *testing.T) {
	engine := testutil.NewTestEngine(t)
	token := testutil.RegisterAndLogin(t, engine)

	// price=10000, stock=5
	productID := testutil.CreateProduct(t, engine, token, "Pants", 10000, 5)

	// create order: qty=2 → subtotal=20000, tax=2200(11%), total=22200
	orderBody, _ := json.Marshal(map[string]any{
		"items": []map[string]any{
			{"product_id": productID, "quantity": 2},
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(orderBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	testutil.AssertStatus(t, w, http.StatusCreated)

	body := testutil.ParseJSON(t, w)
	orderID, _ := body["id"].(string)
	if orderID == "" {
		t.Fatal("expected order ID in response")
	}

	// total = subtotal only (tax applied later in Step Functions workflow)
	// price=10000 * qty=2 = 20000
	total, ok := body["total"].(float64)
	if !ok {
		t.Fatalf("total field missing or wrong type in: %s", w.Body.String())
	}
	if int64(total) != 20000 {
		t.Errorf("total=%d, want=20000", int64(total))
	}

	// get order — status=PENDING, items match
	req = httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+orderID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	testutil.AssertStatus(t, w, http.StatusOK)
	testutil.AssertJSONField(t, w, "status", "PENDING")
}

// Flow 3: order with insufficient stock returns 422
func TestFlow_InsufficientStock(t *testing.T) {
	engine := testutil.NewTestEngine(t)
	token := testutil.RegisterAndLogin(t, engine)

	productID := testutil.CreateProduct(t, engine, token, "Hat", 5000, 1)

	orderBody, _ := json.Marshal(map[string]any{
		"items": []map[string]any{
			{"product_id": productID, "quantity": 99},
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(orderBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	testutil.AssertStatus(t, w, http.StatusUnprocessableEntity)
}

// Flow 4: error handling
func TestFlow_ErrorHandling(t *testing.T) {
	engine := testutil.NewTestEngine(t)

	t.Run("no token returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		testutil.AssertStatus(t, w, http.StatusUnauthorized)
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
		req.Header.Set("Authorization", "Bearer bad-token")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		testutil.AssertStatus(t, w, http.StatusUnauthorized)
	})

	t.Run("duplicate email returns 409", func(t *testing.T) {
		regBody, _ := json.Marshal(map[string]any{
			"name": "User", "email": "dup@example.com",
			"password": "secret123", "phone": "08123456789",
		})
		for i := 0; i < 2; i++ {
			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(regBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			if i == 1 {
				testutil.AssertStatus(t, w, http.StatusConflict)
			}
		}
	})

	t.Run("create product with missing name returns 400", func(t *testing.T) {
		token := testutil.RegisterAndLogin(t, engine)
		body, _ := json.Marshal(map[string]any{"price": 1000, "stock": 5})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		testutil.AssertStatus(t, w, http.StatusBadRequest)
	})
}

// Flow 5: health check (no auth)
func TestFlow_Health(t *testing.T) {
	engine := testutil.NewTestEngine(t)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	testutil.AssertStatus(t, w, http.StatusOK)
}
