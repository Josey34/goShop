package order

import (
	"context"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	"github.com/Josey34/goshop/repository"
)

type ListOrderUseCase struct {
	orderRepo repository.OrderRepository
}

func NewListProductUseCase(repo repository.OrderRepository) *ListOrderUseCase {
	return &ListOrderUseCase{orderRepo: repo}
}

func (uc *ListOrderUseCase) Execute(ctx context.Context, customerID string, pagination valueobject.Pagination) ([]*entity.Order, error) {
	order, err := uc.orderRepo.FindByCustomer(ctx, customerID, pagination)

	if err != nil {
		return nil, err
	}

	return order, nil
}
