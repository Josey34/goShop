package response

type OrderItemResponse struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	Quantity  int    `json:"quantity"`
	LineTotal int64  `json:"line_total"`
}

type OrderResponse struct {
	ID         string              `json:"id"`
	CustomerID string              `json:"customer_id"`
	Items      []OrderItemResponse `json:"items"`
	Total      int64               `json:"total"`
	Status     string              `json:"status"`
	CreatedAt  string              `json:"created_at"`
}
