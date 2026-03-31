package factory

import (
	"fmt"
	"time"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
)

type CreateOrderInput struct {
	CustomerID string
	Items      []CreateOrderItemInput
}

type CreateOrderItemInput struct {
	ProductID string
	Name      string
	Price     valueobject.Money
	Quantity  int
}

func NewOrder(id string, input CreateOrderInput) (*entity.Order, error) {
	if input.Items == nil {
		return nil, fmt.Errorf("order must have at least one item")
	}

	if len(input.Items) == 0 {
		return nil, fmt.Errorf("order must have at least one item")
	}

	var items []entity.OrderItem
	for i, item := range input.Items {
		items = append(items, entity.OrderItem{
			ID:        fmt.Sprintf("%s-item-%d", id, i),
			OrderID:   id,
			ProductID: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
		})
	}

	order := &entity.Order{
		ID:         id,
		CustomerID: input.CustomerID,
		Items:      items,
		Status:     valueobject.OrderStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	order.CalculateTotal()
	return order, nil
}
