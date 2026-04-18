package product

import (
	"context"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/repository"
)

type GetProductUseCase struct {
	productRepo repository.ProductRepository
}

func NewGetProductUseCase(repo repository.ProductRepository) *GetProductUseCase {
	return &GetProductUseCase{productRepo: repo}
}

func (uc *GetProductUseCase) Execute(ctx context.Context, id string) (*entity.Product, error) {
	product, err := uc.productRepo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}
