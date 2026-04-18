package auth

import (
	"context"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/errors"
	"github.com/Josey34/goshop/pkg/hasher"
	"github.com/Josey34/goshop/repository"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginUseCase struct {
	customerRepo repository.CustomerRepository
	hasher       hasher.PasswordHasher
}

func NewLoginUseCase(customerRepo repository.CustomerRepository, hasher hasher.PasswordHasher) *LoginUseCase {
	return &LoginUseCase{
		customerRepo: customerRepo,
		hasher:       hasher,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*entity.Customer, error) {
	customer, err := uc.customerRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if err := uc.hasher.Compare(customer.PasswordHash, input.Password); err != nil {
		return nil, errors.NewUnauthorized("invalid credentials")
	}

	return customer, nil
}
