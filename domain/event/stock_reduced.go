package event

import "time"

type StockReduced struct {
	ProductID  string
	Quantity   int
	Remaining  int
	OccurredOn time.Time
}

func NewStockReduced(productID string, quantity, remaining int) StockReduced {
	return StockReduced{
		ProductID:  productID,
		Quantity:   quantity,
		Remaining:  remaining,
		OccurredOn: time.Now(),
	}
}

func (e StockReduced) EventName() string {
	return "stock.reduced"
}

func (e StockReduced) OccurredAt() time.Time {
	return e.OccurredOn
}
