package main

type ProcessPaymentEvent struct {
	OrderID string `json:"order_id"`
	Amount  int64  `json:"amount"`
}

type ProcessPaymentResponse struct {
	OrderID       string `json:"order_id"`
	PaymentStatus string `json:"payment_status"`
	Message       string `json:"message"`
}
