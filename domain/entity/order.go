package entity

import (
	"time"

	"github.com/Josey34/goshop/domain/errors"
	"github.com/Josey34/goshop/domain/valueobject"
)

type Order struct {
	ID         string
	CustomerID string
	Items      []OrderItem
	Total      valueobject.Money
	Status     valueobject.OrderStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (o *Order) CalculateTotal() {
	var total valueobject.Money

	for _, item := range o.Items {
		total = total.Add(item.LineTotal())
	}

	o.Total = total
}

func (o *Order) TransitionTo(next valueobject.OrderStatus) error {
	if !o.Status.CanTransitionTo(next) {
		return errors.NewInvalidTransition(string(o.Status), string(next))
	}

	o.Status = next
	o.UpdatedAt = time.Now()

	return nil
}
