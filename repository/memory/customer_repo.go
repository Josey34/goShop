package memory

import (
	"context"
	"sync"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/errors"
)

type CustomerRepo struct {
	mu        sync.RWMutex
	customers map[string]*entity.Customer
}

func NewCustomerRepo() *CustomerRepo {
	return &CustomerRepo{customers: make(map[string]*entity.Customer)}
}

func (r *CustomerRepo) Create(ctx context.Context, customer *entity.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.customers[customer.ID] = customer
	return nil
}

func (r *CustomerRepo) FindByID(ctx context.Context, id string) (*entity.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.customers[id]
	if !ok {
		return nil, errors.NewNotFound("customer", id)
	}

	return p, nil
}

func (r *CustomerRepo) FindByEmail(ctx context.Context, email string) (*entity.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, c := range r.customers {
		if c.Email.Value() == email {
			return c, nil
		}
	}
	return nil, errors.NewNotFound("customer", email)
}

func (r *CustomerRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, c := range r.customers {
		if c.Email.Value() == email {
			return false, nil
		}
	}
	return false, nil
}
