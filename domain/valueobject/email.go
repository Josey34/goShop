package valueobject

import (
	"errors"
	"strings"
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	if value == "" {
		return Email{}, errors.New("email cannot be empty")
	}

	if !strings.Contains(value, "@") {
		return Email{}, errors.New("Invalid email format")
	}

	return Email{value: value}, nil
}

func (e Email) Value() string {
	return e.value
}
