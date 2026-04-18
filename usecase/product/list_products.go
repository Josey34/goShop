package product

import (
	"context"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	"github.com/Josey34/goshop/repository"
)

type ListProductsUseCase struct {
	productRepo repository.ProductRepository
}

func NewListProductUseCase(repo repository.ProductRepository) *ListProductsUseCase {
	return &ListProductsUseCase{productRepo: repo}
}

func (uc *ListProductsUseCase) Execute(ctx context.Context, pagination valueobject.Pagination) ([]*entity.Product, error) {
	product, err := uc.productRepo.FindAll(ctx, pagination)

	if err != nil {
		return nil, err
	}

	return product, nil
}
