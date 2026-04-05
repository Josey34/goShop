package repository

import (
	"context"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	FindByID(ctx context.Context, id string) (*entity.Order, error)
	FindByCustomer(ctx context.Context, customerID string, pagination valueobject.Pagination) ([]*entity.Order, error)
	UpdateStatus(ctx context.Context, id string, status valueobject.OrderStatus) error
}
