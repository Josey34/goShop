package memory

import (
	"context"
	"sync"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/errors"
	"github.com/Josey34/goshop/domain/valueobject"
)

type ProductRepo struct {
	mu       sync.RWMutex
	products map[string]*entity.Product
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{
		products: make(map[string]*entity.Product),
	}
}

func (r *ProductRepo) Create(ctx context.Context, product *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.products[product.ID] = product
	return nil
}

func (r *ProductRepo) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.products[id]
	if !ok {
		return nil, errors.NewNotFound("product", id)
	}

	return p, nil
}

func (r *ProductRepo) FindAll(ctx context.Context, pagination valueobject.Pagination) ([]*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	all := make([]*entity.Product, 0, len(r.products))
	for _, p := range r.products {
		all = append(all, p)
	}

	offset := pagination.Offset()
	if offset >= len(all) {
		return []*entity.Product{}, nil
	}
	end := offset + pagination.Limit
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], nil
}

func (r *ProductRepo) Update(ctx context.Context, product *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[product.ID]; !ok {
		return errors.NewNotFound("product", product.ID)
	}

	r.products[product.ID] = product
	return nil
}

func (r *ProductRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.products, id)
	return nil
}
