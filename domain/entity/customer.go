package entity

import (
	"time"

	"github.com/Josey34/goshop/domain/valueobject"
)

type Customer struct {
	ID           string
	Name         string
	Email        valueobject.Email
	Phone        valueobject.Phone
	Address      valueobject.Address
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
