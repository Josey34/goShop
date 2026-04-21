package request

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}
