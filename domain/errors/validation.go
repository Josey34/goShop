package errors

func NewValidation(message string, fieldErrors map[string]string) *DomainError {
	return &DomainError{
		Code:    CodeValidation,
		Message: message + " validation failed",
		Fields:  fieldErrors,
	}
}
