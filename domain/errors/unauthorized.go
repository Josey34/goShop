package errors

func NewUnauthorized(message string) *DomainError {
	return &DomainError{
		Code:    CodeUnauthorized,
		Message: message,
	}
}
