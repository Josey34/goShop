package main

type CalculateTotalEvent struct {
	OrderID string `json:"order_id"`
}

type CalculateTotalResponse struct {
	OrderID  string `json:"order_id"`
	Subtotal int64  `json:"subtotal"`
	Tax      int64  `json:"tax"`
	Total    int64  `json:"total"`
}
