package event

import "time"

type OrderCreated struct {
	OrderID    string
	CustomerID string
	Total      int64
	OccurredOn time.Time
}

func NewOrderCreated(orderID, customerID string, total int64) OrderCreated {
	return OrderCreated{
		OrderID:    orderID,
		CustomerID: customerID,
		Total:      total,
		OccurredOn: time.Now(),
	}
}

func (e OrderCreated) EventName() string {
	return "order.created"
}

func (e OrderCreated) OccurredAt() time.Time {
	return e.OccurredOn
}
