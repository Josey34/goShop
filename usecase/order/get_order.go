package order

import (
	"context"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/repository"
)

type GetOrderUseCase struct {
	orderRepo repository.OrderRepository
}

func NewGetOrderUseCase(repo repository.OrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{orderRepo: repo}
}

func (uc *GetOrderUseCase) Execute(ctx context.Context, id string) (*entity.Order, error) {
	product, err := uc.orderRepo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}
