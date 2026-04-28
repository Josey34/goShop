package main

type ValidateOrderEvent struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
}

type ValidateOrderResponse struct {
	OrderID string   `json:"order_id"`
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors"`
}
