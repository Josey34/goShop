package request

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       int64  `json:"price" binding:"required,gt=0"`
	Stock       int    `json:"stock" binding:"gte=0"`
}
