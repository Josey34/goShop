package errors

func NewNotFound(message string, id string) *DomainError {
	return &DomainError{
		Code:    CodeNotFound,
		Message: message + " with id " + id + " was not found",
	}
}
