package product

import (
	"context"
	"time"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	"github.com/Josey34/goshop/repository"
)

type UpdateProductInput struct {
	ID          string
	Name        string
	Description string
	Price       int64
	Stock       int
}

type UpdateProductUseCase struct {
	productRepo repository.ProductRepository
}

func NewUpdateProductUseCase(repo repository.ProductRepository) *UpdateProductUseCase {
	return &UpdateProductUseCase{productRepo: repo}
}

func (uc *UpdateProductUseCase) Execute(ctx context.Context, input UpdateProductInput) (*entity.Product, error) {
	productFound, err := uc.productRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	product := entity.Product{
		ID:          productFound.ID,
		Name:        input.Name,
		Description: input.Description,
		Price:       valueobject.NewMoney(input.Price),
		Stock:       input.Stock,
		CreatedAt:   productFound.CreatedAt,
		UpdatedAt:   now,
	}

	if err := uc.productRepo.Update(ctx, &product); err != nil {
		return nil, err
	}

	return &product, nil
}
