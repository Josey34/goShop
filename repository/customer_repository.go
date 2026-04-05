package repository

import (
	"context"

	"github.com/Josey34/goshop/domain/entity"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *entity.Customer) error
	FindByID(ctx context.Context, id string) (*entity.Customer, error)
	FindByEmail(ctx context.Context, email string) (*entity.Customer, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
