package memory

import (
	"context"
	"sync"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/errors"
	"github.com/Josey34/goshop/domain/valueobject"
)

type OrderRepo struct {
	mu     sync.RWMutex
	orders map[string]*entity.Order
}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{
		orders: make(map[string]*entity.Order),
	}
}

func (r *OrderRepo) Create(ctx context.Context, order *entity.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

func (r *OrderRepo) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[id]
	if !ok {
		return nil, errors.NewNotFound("order", id)
	}
	return o, nil
}

func (r *OrderRepo) FindByCustomer(ctx context.Context, customerID string, pagination valueobject.Pagination) ([]*entity.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matched []*entity.Order
	for _, o := range r.orders {
		if o.CustomerID == customerID {
			matched = append(matched, o)
		}
	}

	offset := pagination.Offset()
	if offset >= len(matched) {
		return []*entity.Order{}, nil
	}
	end := offset + pagination.Limit
	if end > len(matched) {
		end = len(matched)
	}
	return matched[offset:end], nil
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, id string, status valueobject.OrderStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	o, ok := r.orders[id]
	if !ok {
		return errors.NewNotFound("order", id)
	}

	if err := o.TransitionTo(status); err != nil {
		return err
	}

	r.orders[id] = o
	return nil
}
