package valueobject

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusConfirmed  OrderStatus = "CONFIRMED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusShipped    OrderStatus = "SHIPPED"
	OrderStatusDelivered  OrderStatus = "DELIVERED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
)

var validTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusPending:    {OrderStatusConfirmed, OrderStatusCancelled},
	OrderStatusConfirmed:  {OrderStatusProcessing, OrderStatusCancelled},
	OrderStatusProcessing: {OrderStatusShipped},
	OrderStatusShipped:    {OrderStatusDelivered},
	OrderStatusDelivered:  {},
	OrderStatusCancelled:  {},
}

func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	for _, status := range validTransitions[s] {
		if status == next {
			return true
		}
	}
	return false
}
