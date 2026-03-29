package errors

import "fmt"

type ErrorCode string

const (
	CodeNotFound          ErrorCode = "NOT_FOUND"
	CodeValidation        ErrorCode = "VALIDATION"
	CodeConflict          ErrorCode = "CONFLICT"
	CodeUnauthorized      ErrorCode = "UNAUTHORIZED"
	CodeInsufficientStock ErrorCode = "INSUFFICIENT_STOCK"
	CodeInvalidTransition ErrorCode = "INVALID_TRANSITION"
)

type DomainError struct {
	Code    ErrorCode
	Message string
	Fields  map[string]string
}

func (e *DomainError) Error() string {
	if e.Fields != nil {
		return fmt.Sprintf("Error %s: %s - Fields: %v", e.Code, e.Message, e.Fields)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func IsNotFound(err error) bool {
	if err, ok := err.(*DomainError); ok {
		return err.Code == CodeNotFound
	}
	return false
}

func IsValidation(err error) bool {
	if err, ok := err.(*DomainError); ok {
		return err.Code == CodeValidation
	}
	return false
}

func IsConflict(err error) bool {
	if err, ok := err.(*DomainError); ok {
		return err.Code == CodeConflict
	}
	return false
}

func IsUnauthorized(err error) bool {
	if err, ok := err.(*DomainError); ok {
		return err.Code == CodeUnauthorized
	}
	return false
}

func IsInsufficientStock(err error) bool {
	if err, ok := err.(*DomainError); ok {
		return err.Code == CodeInsufficientStock
	}
	return false
}

func IsInvalidTransition(err error) bool {
	if err, ok := err.(*DomainError); ok {
		return err.Code == CodeInvalidTransition
	}
	return false
}
