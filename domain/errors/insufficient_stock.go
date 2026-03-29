package errors

import "fmt"

func NewInsufficientStock(productID string, requested int, available int) *DomainError {
	return &DomainError{
		Code:    CodeInsufficientStock,
		Message: fmt.Sprintf("Insufficient stock for product %s. Requested: %d, Available: %d", productID, requested, available),
	}
}
