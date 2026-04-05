package repository

import (
	"context"

	"github.com/Josey34/goshop/domain/entity"
)

type OrderQueueRepository interface {
	Enqueue(ctx context.Context, order *entity.Order) error
}
