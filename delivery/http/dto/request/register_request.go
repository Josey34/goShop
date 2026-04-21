package request

type RegisterRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	Phone      string `json:"phone" binding:"required"`
	Street     string `json:"street"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`
}
