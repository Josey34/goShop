package auth

import (
	"context"

	"github.com/Josey34/goshop/domain/errors"
	"github.com/Josey34/goshop/domain/factory"
	"github.com/Josey34/goshop/pkg/hasher"
	"github.com/Josey34/goshop/pkg/idgen"
	"github.com/Josey34/goshop/repository"
)

type RegisterInput struct {
	Name       string
	Email      string
	Phone      string
	Street     string
	City       string
	Province   string
	PostalCode string
	Password   string
}

type RegisiterUseCase struct {
	customerRepo repository.CustomerRepository
	hasher       hasher.PasswordHasher
	idGen        idgen.IDGenerator
}

func NewRegisterUseCase(customerRepo repository.CustomerRepository, hasher hasher.PasswordHasher, idGen idgen.IDGenerator) *RegisiterUseCase {
	return &RegisiterUseCase{
		customerRepo: customerRepo,
		hasher:       hasher,
		idGen:        idGen,
	}
}

func (uc *RegisiterUseCase) Execute(ctx context.Context, input RegisterInput) error {
	found, err := uc.customerRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return err
	}

	if found {
		return errors.NewConflict("customer", input.Email)
	}

	hashedPassword, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	customer, err := factory.NewCustomer(uc.idGen.Generate(), factory.CreateCustomerInput{
		Name:       input.Name,
		Email:      input.Email,
		Phone:      input.Phone,
		Street:     input.Street,
		City:       input.City,
		Province:   input.Province,
		PostalCode: input.PostalCode,
	}, hashedPassword)
	if err != nil {
		return err
	}

	if err := uc.customerRepo.Create(ctx, customer); err != nil {
		return err
	}

	return nil
}
