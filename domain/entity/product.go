package entity

import (
	"time"

	"github.com/Josey34/goshop/domain/errors"
	"github.com/Josey34/goshop/domain/valueobject"
)

type Product struct {
	ID          string
	Name        string
	Description string
	Price       valueobject.Money
	Stock       int
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Product) ReduceStock(quantity int) error {
	if p.Stock < quantity {
		return errors.NewInsufficientStock(p.ID, quantity, p.Stock)
	}

	p.Stock -= quantity
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Product) UpdatePrice(price valueobject.Money) error {
	if price.IsZero() || price.IsNegative() {
		return errors.NewValidation("price", map[string]string{"price": "must be greater than zero"})
	}

	p.Price = price
	p.UpdatedAt = time.Now()

	return nil
}
