package event

import "time"

type OrderStatusChanged struct {
	OrderID    string
	From       string
	To         string
	OccurredOn time.Time
}

func NewOrderStatusChanged(orderID, from, to string) OrderStatusChanged {
	return OrderStatusChanged{
		OrderID:    orderID,
		From:       from,
		To:         to,
		OccurredOn: time.Now(),
	}
}

func (e OrderStatusChanged) EventName() string {
	return "order.status_changed"
}

func (e OrderStatusChanged) OccurredAt() time.Time {
	return e.OccurredOn
}
