package main

type SendNotificationEvent struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	Message    string `json:"message"`
}

type SendNotificationResponse struct {
	OrderID string `json:"order_id"`
	Sent    bool   `json:"sent"`
}
