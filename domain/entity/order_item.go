package entity

import "github.com/Josey34/goshop/domain/valueobject"

type OrderItem struct {
	ID        string
	OrderID   string
	ProductID string
	Name      string
	Price     valueobject.Money
	Quantity  int
}

func (i OrderItem) LineTotal() valueobject.Money {
	return i.Price.Multiply(i.Quantity)
}
