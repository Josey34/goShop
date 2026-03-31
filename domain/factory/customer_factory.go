package factory

import (
	"time"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
)

type CreateCustomerInput struct {
	Name       string
	Email      string
	Phone      string
	Street     string
	City       string
	Province   string
	PostalCode string
	Password   string
}

func NewCustomer(id string, input CreateCustomerInput, hashedPassword string) (*entity.Customer, error) {
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	phone, err := valueobject.NewPhone(input.Phone)
	if err != nil {
		return nil, err
	}

	address, err := valueobject.NewAddress(input.Street, input.City, input.Province, input.PostalCode)
	if err != nil {
		return nil, err
	}

	return &entity.Customer{
		ID:           id,
		Name:         input.Name,
		Email:        email,
		Phone:        phone,
		Address:      address,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
