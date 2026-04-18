package product

import (
	"context"
	"time"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/errors"
	"github.com/Josey34/goshop/domain/valueobject"
	"github.com/Josey34/goshop/pkg/idgen"
	"github.com/Josey34/goshop/repository"
)

type CreateProductInput struct {
	Name        string
	Description string
	Price       int64
	Stock       int
}

type CreateProductUseCase struct {
	productRepo repository.ProductRepository
	idGen       idgen.IDGenerator
}

func NewCreateProductUseCase(repo repository.ProductRepository, idGen idgen.IDGenerator) *CreateProductUseCase {
	return &CreateProductUseCase{
		productRepo: repo,
		idGen:       idGen,
	}
}

func (uc *CreateProductUseCase) Execute(ctx context.Context, input CreateProductInput) (*entity.Product, error) {
	if input.Name == "" {
		return nil, errors.NewValidation("name", map[string]string{"name": "Name is required"})
	}

	now := time.Now()
	product := entity.Product{
		ID:          uc.idGen.Generate(),
		Name:        input.Name,
		Description: input.Description,
		Price:       valueobject.NewMoney(input.Price),
		Stock:       input.Stock,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.productRepo.Create(ctx, &product); err != nil {
		return nil, err
	}

	return &product, nil
}
