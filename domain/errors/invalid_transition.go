package errors

func NewInvalidTransition(from string, to string) *DomainError {
	return &DomainError{
		Code:    CodeInvalidTransition,
		Message: "Invalid transition from " + from + " to " + to,
	}
}
