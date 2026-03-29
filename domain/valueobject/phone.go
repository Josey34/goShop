package valueobject

import (
	"errors"
	"strings"
)

type Phone struct {
	value string
}

func NewPhone(value string) (Phone, error) {
	if value == "" {
		return Phone{}, errors.New("phone is required")
	}

	if !strings.HasPrefix(value, "08") && !strings.HasPrefix(value, "+62") {
		return Phone{}, errors.New("phone must start with '08' or '+62'")
	}

	return Phone{value: value}, nil
}

func (p Phone) Value() string {
	return p.value
}
