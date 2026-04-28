package main

type FulfillOrderEvent struct {
	OrderID string `json:"order_id"`
}

type FullfillOrderResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
