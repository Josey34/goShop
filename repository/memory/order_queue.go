package memory

import (
	"context"
	"fmt"

	"github.com/Josey34/goshop/domain/entity"
)

type OrderQueue struct{}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{}
}

func (q *OrderQueue) Enqueue(ctx context.Context, order *entity.Order) error {
	fmt.Printf("order queued: %s\n", order.ID)
	return nil
}
