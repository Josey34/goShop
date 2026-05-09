package order_test

import (
	"context"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
)

type mockOrderRepo struct {
	created    *entity.Order
	findResult *entity.Order
	findErr    error
	listResult []*entity.Order
	listErr    error
	updateErr  error
}

func (m *mockOrderRepo) Create(ctx context.Context, o *entity.Order) error {
	m.created = o
	return nil
}
func (m *mockOrderRepo) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	return m.findResult, m.findErr
}
func (m *mockOrderRepo) FindByCustomer(ctx context.Context, customerID string, pagination valueobject.Pagination) ([]*entity.Order, error) {
	return m.listResult, m.listErr
}
func (m *mockOrderRepo) UpdateStatus(ctx context.Context, id string, status valueobject.OrderStatus) error {
	return m.updateErr
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

func makeOrder(id string, status valueobject.OrderStatus) *entity.Order {
	return &entity.Order{ID: id, CustomerID: "c1", Status: status}
}
