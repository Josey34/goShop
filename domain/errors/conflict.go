package errors

func NewConflict(message string, value string) *DomainError {
	return &DomainError{
		Code:    CodeConflict,
		Message: message + " already exists with value: " + value,
	}
}
