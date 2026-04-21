package response

type AuthResponse struct {
	Token      string `json:"token"`
	CustomerID string `json:"customer_id"`
	Email      string `json:"email"`
}
