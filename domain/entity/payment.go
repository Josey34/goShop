package entity

import (
	"time"

	"github.com/Josey34/goshop/domain/valueobject"
)

type Payment struct {
	ID        string
	OrderID   string
	Amount    valueobject.Money
	Status    valueobject.PaymentStatus
	Method    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
