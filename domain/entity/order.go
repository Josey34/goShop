package entity

import (
	"time"

	"github.com/Josey34/goshop/domain/errors"
	"github.com/Josey34/goshop/domain/event"
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
	events     []event.DomainEvent
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

func (o *Order) AddEvent(e event.DomainEvent) {
	o.events = append(o.events, e)
}

func (o *Order) PullEvents() []event.DomainEvent {
	events := o.events

	o.events = nil

	return events
}
