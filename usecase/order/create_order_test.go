package order_test

import (
	"context"
	"testing"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	"github.com/Josey34/goshop/domain/factory"
	ucorder "github.com/Josey34/goshop/usecase/order"
)

// --- mocks ---

type mockOrderRepo struct {
	created *entity.Order
}

func (m *mockOrderRepo) Create(ctx context.Context, o *entity.Order) error {
	m.created = o
	return nil
}
func (m *mockOrderRepo) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	return nil, nil
}
func (m *mockOrderRepo) FindByCustomer(ctx context.Context, customerID string, pagination valueobject.Pagination) ([]*entity.Order, error) {
	return nil, nil
}
func (m *mockOrderRepo) UpdateStatus(ctx context.Context, id string, status valueobject.OrderStatus) error {
	return nil
}

type mockProductRepo struct {
	products map[string]*entity.Product
}

func (m *mockProductRepo) Create(ctx context.Context, p *entity.Product) error { return nil }
func (m *mockProductRepo) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	p, ok := m.products[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}
func (m *mockProductRepo) FindAll(ctx context.Context, pagination valueobject.Pagination) ([]*entity.Product, error) {
	return nil, nil
}
func (m *mockProductRepo) Update(ctx context.Context, p *entity.Product) error { return nil }
func (m *mockProductRepo) Delete(ctx context.Context, id string) error         { return nil }

type mockIDGen struct{ id string }

func (m *mockIDGen) Generate() string { return m.id }

// --- tests ---

func TestCreateOrderUseCase_Execute(t *testing.T) {
	t.Run("valid order reduces stock", func(t *testing.T) {
		productRepo := &mockProductRepo{
			products: map[string]*entity.Product{
				"p1": {ID: "p1", Name: "Shirt", Price: valueobject.NewMoney(1000), Stock: 10},
			},
		}
		orderRepo := &mockOrderRepo{}
		idGen := &mockIDGen{id: "order-1"}
		uc := ucorder.NewCreateOrderUseCase(orderRepo, productRepo, idGen)

		order, err := uc.Execute(context.Background(), ucorder.CreateOrderInput{
			CustomerID: "c1",
			Items:      []factory.CreateOrderItemInput{{ProductID: "p1", Quantity: 2}},
		})

		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if order.ID != "order-1" {
			t.Errorf("got id=%s, want=order-1", order.ID)
		}
		if productRepo.products["p1"].Stock != 8 {
			t.Errorf("got stock=%d, want=8", productRepo.products["p1"].Stock)
		}
		if order.Total.Value() != 2000 {
			t.Errorf("got total=%d, want=2000", order.Total.Value())
		}
	})

	t.Run("insufficient stock returns error", func(t *testing.T) {
		productRepo := &mockProductRepo{
			products: map[string]*entity.Product{
				"p1": {ID: "p1", Name: "Shirt", Price: valueobject.NewMoney(1000), Stock: 1},
			},
		}
		uc := ucorder.NewCreateOrderUseCase(&mockOrderRepo{}, productRepo, &mockIDGen{id: "x"})

		_, err := uc.Execute(context.Background(), ucorder.CreateOrderInput{
			CustomerID: "c1",
			Items:      []factory.CreateOrderItemInput{{ProductID: "p1", Quantity: 5}},
		})

		if err == nil {
			t.Error("expected error for insufficient stock")
		}
	})
}
