package product

import (
	"context"

	"github.com/Josey34/goshop/repository"
)

type DeleteProductUseCase struct {
	productRepo repository.ProductRepository
}

func NewDeleteProductUseCase(repo repository.ProductRepository) *DeleteProductUseCase {
	return &DeleteProductUseCase{productRepo: repo}
}

func (uc *DeleteProductUseCase) Execute(ctx context.Context, id string) error {
	_, err := uc.productRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return uc.productRepo.Delete(ctx, id)
}
