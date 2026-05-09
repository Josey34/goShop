package product_test

import (
	"context"
	"time"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
)

type mockProductRepo struct {
	createErr  error
	findResult *entity.Product
	findErr    error
	listResult []*entity.Product
	listErr    error
	updateErr  error
	deleteErr  error
}

func (m *mockProductRepo) Create(ctx context.Context, p *entity.Product) error { return m.createErr }
func (m *mockProductRepo) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	return m.findResult, m.findErr
}
func (m *mockProductRepo) FindAll(ctx context.Context, pagination valueobject.Pagination) ([]*entity.Product, error) {
	return m.listResult, m.listErr
}
func (m *mockProductRepo) Update(ctx context.Context, p *entity.Product) error { return m.updateErr }
func (m *mockProductRepo) Delete(ctx context.Context, id string) error         { return m.deleteErr }

type mockIDGen struct{ id string }

func (m *mockIDGen) Generate() string { return m.id }

func newProduct(id, name string, price int64) *entity.Product {
	return &entity.Product{
		ID:        id,
		Name:      name,
		Price:     valueobject.NewMoney(price),
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
