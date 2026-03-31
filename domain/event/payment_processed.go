package event

import "time"

type PaymentProcessed struct {
	OrderID    string
	PaymentID  string
	Amount     int64
	Status     string
	OccurredOn time.Time
}

func NewPaymentProcessed(orderID, paymentID string, amount int64, status string) PaymentProcessed {
	return PaymentProcessed{
		OrderID:    orderID,
		PaymentID:  paymentID,
		Amount:     amount,
		Status:     status,
		OccurredOn: time.Now(),
	}
}

func (e PaymentProcessed) EventName() string {
	return "payment.processed"
}

func (e PaymentProcessed) OccurredAt() time.Time {
	return e.OccurredOn
}
